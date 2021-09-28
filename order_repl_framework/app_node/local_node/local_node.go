package local_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"sync/atomic"
	"time"
)

const (
	batchInterval = 100 * time.Microsecond
)

type LocalNode struct {
	app frame.Application

	id      uint32
	color   uint32
	ctr     uint32
	helpCtr uint32

	comReqCh chan *frame.CommitRequest
	prepCh   chan *nodepb.Prep
	oReqCh   chan *sequencerpb.OrderRequest
	oRspCh   chan *sequencerpb.OrderResponse

	prepM *prepManager
}

func NewLocalNode(id, color uint32, app frame.Application, onPrep frame.OnPrepFunc, onOrder frame.OnOrderFunc) *LocalNode {
	ln := new(LocalNode)
	ln.id = id
	ln.color = color
	ln.app = app
	ln.comReqCh = make(chan *frame.CommitRequest, 1024)
	ln.prepCh = make(chan *nodepb.Prep, 1024)
	ln.oReqCh = make(chan *sequencerpb.OrderRequest, 1024)
	ln.oRspCh = make(chan *sequencerpb.OrderResponse, 1024)
	ln.prepM = newPrepManager(app)

	go ln.commit()
	go ln.handleCommitRequests()
	go ln.batchPrepMsg(batchInterval, onPrep)
	go ln.batchOrderRequest(batchInterval, onOrder)
	return ln
}

func (ln *LocalNode) PutComReq(comReq *frame.CommitRequest) {
	ln.comReqCh <- comReq
}

func (ln *LocalNode) Prepare(prepMsg *nodepb.Prep) {
	ln.prepM.prepare(prepMsg, 0)
}

func (ln *LocalNode) PutOrderResponse(oRsp *sequencerpb.OrderResponse) {
	ln.oRspCh <- oRsp
}

func (ln *LocalNode) commit() {
	for oRsp := range ln.oRspCh {
		isCoor := uint32(oRsp.Lsn>>32) == ln.id
		for i := uint64(0); i < uint64(oRsp.NumOfRecords); i++ {
			ln.prepM.waitForPrep(oRsp.Lsn + i)
			if err := ln.app.Commit(oRsp.Lsn+i, oRsp.Color, oRsp.Gsn+i, isCoor); err != nil {
				logrus.Fatalln(err)
			}
		}
	}
}

func (ln *LocalNode) handleCommitRequests() {
	for comReq := range ln.comReqCh {
		localToken := ln.getNewToken()
		ln.prepM.prepare(
			&nodepb.Prep{LocalToken: localToken, Color: comReq.Color, Contents: []string{comReq.Content}}, comReq.FindToken)
		ln.prepCh <- &nodepb.Prep{LocalToken: localToken, Color: comReq.Color, Contents: []string{comReq.Content}}
		ln.oReqCh <- &sequencerpb.OrderRequest{Lsn: localToken, NumOfRecords: 1, Color: comReq.Color, OriginColor: ln.color}
	}
}

func (ln *LocalNode) batchPrepMsg(interval time.Duration, onPrep frame.OnPrepFunc) {
	newBatch := false
	var send <-chan time.Time
	var current *nodepb.Prep

	for {
		select {
		case prepMsg := <-ln.prepCh:
			if newBatch {
				send = time.After(interval)
				current = prepMsg
				newBatch = false
				continue
			}
			current.Contents = append(current.Contents, prepMsg.Contents...)

		case <-send:
			onPrep(current)
			newBatch = true
		}
	}
}

func (ln *LocalNode) batchOrderRequest(interval time.Duration, onOrder frame.OnOrderFunc) {
	newBatch := false
	var send <-chan time.Time
	var current *sequencerpb.OrderRequest

	for {
		select {
		case oReq := <-ln.oReqCh:
			if newBatch {
				send = time.After(interval)
				current = oReq
				newBatch = false
				continue
			}
			current.NumOfRecords += 1
		case <-send:
			onOrder(current)
			newBatch = true
		}
	}
}

func (ln *LocalNode) getNewToken() uint64 {
	return (uint64(ln.id) << 32) + uint64(atomic.AddUint32(&ln.ctr, 1))
}
