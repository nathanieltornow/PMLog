package local_node

import (
	"fmt"
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

type LocalNode struct {
	mu sync.Mutex

	app frame.Application

	id          uint32
	peers       uint32
	originColor uint32
	ctr         uint32
	helpCtr     uint32

	comReqMu sync.Mutex

	prepCh chan *nodepb.Prep
	oReqCh chan *sequencerpb.OrderRequest
	oRspCh chan *sequencerpb.OrderResponse

	ackCh chan *nodepb.Acknowledgement

	recentlyCommitted map[uint64]chan bool

	prepM *prepManager
}

func NewLocalNode(id, peers, color uint32, app frame.Application, onPrep frame.OnPrepFunc, onOrder frame.OnOrderFunc, onComm frame.OnCommFunc, interval time.Duration) *LocalNode {
	ln := new(LocalNode)
	ln.id = id
	ln.peers = peers
	ln.originColor = color
	ln.app = app
	ln.prepCh = make(chan *nodepb.Prep, 1024)
	ln.ackCh = make(chan *nodepb.Acknowledgement, 1024)
	ln.oReqCh = make(chan *sequencerpb.OrderRequest, 1024)
	ln.oRspCh = make(chan *sequencerpb.OrderResponse, 1024)
	ln.prepM = newPrepManager(app)
	ln.recentlyCommitted = make(map[uint64]chan bool)

	go ln.handleOrderResponses(onComm)
	go ln.batchPrepMessages(interval, onPrep)
	go ln.batchOrderRequests(interval, onOrder)
	go ln.handleAcknowledgments()
	return ln
}

func (ln *LocalNode) PutComReq(comReq *frame.CommitRequest) uint64 {
	ln.comReqMu.Lock()
	localToken := (uint64(ln.id) << 32) + uint64(ln.ctr)
	ln.ctr++
	ln.prepM.prepare(&nodepb.Prep{LocalToken: localToken, Color: comReq.Color, Contents: []string{comReq.Content}})
	ln.prepCh <- &nodepb.Prep{LocalToken: localToken, Color: comReq.Color, Contents: []string{comReq.Content}}
	ln.oReqCh <- &sequencerpb.OrderRequest{Lsn: localToken, NumOfRecords: 1, Color: comReq.Color, OriginColor: ln.originColor}
	ln.comReqMu.Unlock()
	return localToken
}

func (ln *LocalNode) Prepare(prepMsg *nodepb.Prep) {
	ln.prepM.prepare(prepMsg)
}

func (ln *LocalNode) PutOrderResponse(oRsp *sequencerpb.OrderResponse) {
	ln.oRspCh <- oRsp
}

func (ln *LocalNode) PutAcknowledgment(ack *nodepb.Acknowledgement) {
	ln.ackCh <- ack
}

func (ln *LocalNode) SetPeers(peers uint32) {
	atomic.StoreUint32(&ln.peers, peers)
}

func (ln *LocalNode) handleAcknowledgments() {
	numOfAcks := make(map[uint64]uint32)
	for ack := range ln.ackCh {
		numOfAcks[ack.GlobalToken]++
		if numOfAcks[ack.GlobalToken] == atomic.LoadUint32(&ln.peers) {
			for i := uint64(0); i < uint64(ack.NumOfRecords); i++ {
				ln.waitForCommit(ack.GlobalToken + i)
				if err := ln.app.Acknowledge(ack.LocalToken+i, ack.Color, ack.GlobalToken+i); err != nil {
					logrus.Fatalln(err)
				}
			}
		}
	}
}

func (ln *LocalNode) handleOrderResponses(onComm frame.OnCommFunc) {
	for oRsp := range ln.oRspCh {
		for i := uint64(0); i < uint64(oRsp.NumOfRecords); i++ {
			start := time.Now()
			ln.prepM.waitForPrep(oRsp.Lsn + i)
			fmt.Println(time.Since(start))
			ln.commit(oRsp.Lsn+i, oRsp.Color, oRsp.Gsn+i)
		}
		if uint32(oRsp.Lsn>>32) != ln.id && oRsp.NumOfRecords > 0 {
			onComm(&nodepb.Acknowledgement{Color: oRsp.Color,
				NumOfRecords: oRsp.NumOfRecords,
				LocalToken:   oRsp.Lsn,
				GlobalToken:  oRsp.Gsn})
		}
	}
}

func (ln *LocalNode) batchPrepMessages(interval time.Duration, onPrep frame.OnPrepFunc) {
	newBatch := true
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

func (ln *LocalNode) batchOrderRequests(interval time.Duration, onOrder frame.OnOrderFunc) {
	newBatch := true
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

func (ln *LocalNode) commit(lsn uint64, color uint32, gsn uint64) {
	if err := ln.app.Commit(lsn, color, gsn); err != nil {
		logrus.Fatalln(err)
	}
	ln.mu.Lock()
	ch, ok := ln.recentlyCommitted[gsn]
	if !ok {
		ch = make(chan bool, 1)
		ln.recentlyCommitted[gsn] = ch
	}
	ln.mu.Unlock()
	ch <- true
}

func (ln *LocalNode) waitForCommit(gsn uint64) {
	ln.mu.Lock()
	ch, ok := ln.recentlyCommitted[gsn]
	if !ok {
		ch = make(chan bool, 1)
		ln.recentlyCommitted[gsn] = ch
	}
	ln.mu.Unlock()
	<-ch
}
