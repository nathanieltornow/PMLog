package color_service

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

const (
	prepInterval  = 100 * time.Microsecond
	orderInterval = 100 * time.Microsecond
)

type ColorService struct {
	nodeID      uint32
	color       uint32
	token       uint32
	originColor uint32

	app frame.Application

	prepCh chan *nodepb.Prep
	oReqCh chan *sequencerpb.OrderRequest
	oRspCh chan *sequencerpb.OrderResponse

	prepMan *prepManager
}

func NewColorService(nodeID, color, originColor uint32, app frame.Application, onPrepB frame.OnPrepBatch, onOrderB frame.OnOrderBatch) *ColorService {
	cs := new(ColorService)
	cs.app = app
	cs.nodeID = nodeID
	cs.color = color
	cs.originColor = originColor
	cs.prepMan = newPrepManager(app, color)
	cs.prepCh = make(chan *nodepb.Prep, 1024)
	cs.oReqCh = make(chan *sequencerpb.OrderRequest, 1024)
	cs.oRspCh = make(chan *sequencerpb.OrderResponse, 1025)

	go cs.commit()
	go cs.batchOrderMsg(orderInterval, onOrderB)
	go cs.batchPrepMsg(prepInterval, onPrepB)
	return cs
}

func (cs *ColorService) PrepareCoordinator(content string, findToken uint64) {
	token := atomic.AddUint32(&cs.token, 1)
	localToken := (uint64(cs.nodeID) << 32) + uint64(token)

	cs.prepMan.prepareCoordinator(content, localToken, findToken)
	cs.prepCh <- &nodepb.Prep{LocalToken: localToken, Contents: []string{content}, Color: cs.color}
	oReq := &sequencerpb.OrderRequest{Color: cs.color, NumOfRecords: 1, LocalToken: localToken, OriginColor: cs.originColor}
	cs.oReqCh <- oReq
}

func (cs *ColorService) PutPrepareRequest(prepMsg *nodepb.Prep) {
	cs.prepMan.prepareReplica(prepMsg)
}

func (cs *ColorService) PutOrderResponse(oRsp *sequencerpb.OrderResponse) {
	cs.oRspCh <- oRsp
}

func (cs *ColorService) commit() {
	for oRsp := range cs.oRspCh {
		for i := uint32(0); i < oRsp.NumOfRecords; i++ {
			cs.prepMan.waitForPrep(oRsp.LocalToken + uint64(i))

			isCoor := uint32(oRsp.LocalToken>>32) == cs.nodeID
			err := cs.app.Commit(oRsp.LocalToken+uint64(i), cs.color, oRsp.GlobalToken+uint64(i), isCoor)
			if err != nil {
				logrus.Fatalln(err)
			}
		}
	}
}

func (cs *ColorService) batchPrepMsg(interval time.Duration, onPrepB frame.OnPrepBatch) {
	newBatch := false
	var send <-chan time.Time
	var current *nodepb.Prep

	for {
		select {
		case prepMsg := <-cs.prepCh:
			if newBatch {
				send = time.After(interval)
				current = prepMsg
				newBatch = false
				continue
			}
			current.Contents = append(current.Contents, prepMsg.Contents...)

		case <-send:
			onPrepB(current)
			newBatch = true
		}
	}
}

func (cs *ColorService) batchOrderMsg(interval time.Duration, onOrderB frame.OnOrderBatch) {
	newBatch := false
	var send <-chan time.Time
	var current *sequencerpb.OrderRequest

	for {
		select {
		case oReq := <-cs.oReqCh:
			if newBatch {
				send = time.After(interval)
				current = oReq
				newBatch = false
				continue
			}
			current.NumOfRecords += 1
		case <-send:
			onOrderB(current)
			newBatch = true
		}
	}
}
