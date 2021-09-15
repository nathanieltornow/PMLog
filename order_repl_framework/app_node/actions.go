package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	seqpb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"sync"
)

func (n *Node) broadcastPrepareMsgs() {
	for prepMsg := range n.prepCh {
		n.prepStreamsMu.RLock()
		for _, stream := range n.prepStreams {
			if err := stream.Send(prepMsg); err != nil {
				logrus.Fatalf("failed to send prepmsg: %v", err)
			}
		}
		n.prepStreamsMu.RUnlock()
	}
}

func (n *Node) broadcastCommitMsgs() {
	for comMsg := range n.comCh {
		n.comStreamsMu.RLock()
		for _, stream := range n.comStreams {
			if err := stream.Send(comMsg); err != nil {
				logrus.Fatalf("failed to send commsg: %v", err)
			}
		}
		n.comStreamsMu.RUnlock()
	}
}

func (n *Node) handleAppCommitRequests() {
	comReqCh := make(chan *frame.CommitRequest, 1024)
	go func() {
		err := n.app.MakeCommitRequests(comReqCh)
		if err != nil {
			return
		}
	}()
	for comReq := range comReqCh {
		newToken := n.getNewLocalToken()

		if err := n.app.Prepare(newToken, comReq.Color, comReq.Content, comReq.FindToken); err != nil {
			return
		}
		// put into channel to be broadcasted to other nodes and send orderrequest

		n.prepCh <- &nodepb.Prep{
			LocalToken: newToken,
			Color:      comReq.Color,
			Content:    comReq.Content,
		}
		n.orderReqCh <- &seqpb.OrderRequest{Lsn: newToken, Color: comReq.Color, OriginColor: n.color}
	}
}

func (n *Node) receiveAcks(stream nodepb.Node_GetAcksClient) {
	for {
		ackMsg, err := stream.Recv()
		if err != nil {
			logrus.Fatalln("failed to receive ack")
		}
		n.numOfAcksMu.Lock()
		num, ok := n.numOfAcks[ackMsg.LocalToken]
		if !ok {
			n.numOfAcks[ackMsg.LocalToken] = 0
		}
		n.numOfAcks[ackMsg.LocalToken] += 1
		n.numOfAcksMu.Unlock()
		if (num + 1) == n.numOfPeers {
			comMsg := &nodepb.Com{
				LocalToken:  ackMsg.LocalToken,
				Color:       ackMsg.Color,
				GlobalToken: ackMsg.GlobalToken,
			}
			n.possibleComCh <- comMsg
			n.comCh <- comMsg

			n.numOfAcksMu.Lock()
			delete(n.numOfAcks, ackMsg.LocalToken)
			n.numOfAcksMu.Unlock()
		}

	}
}

func (n *Node) handleOrderResponses() {
	for oRsp := range n.orderRespCh {
		n.waitingORespCh <- oRsp
		id := uint32(oRsp.Lsn)
		if id == n.id {
			continue
		}
		if !n.app.IsPrepared(oRsp.Lsn) {
			logrus.Fatalln("not prepared")
		}
		n.ackChsMu.RLock()
		n.ackChs[id] <- &nodepb.Ack{
			LocalToken:  oRsp.Lsn,
			GlobalToken: oRsp.Gsn,
			Color:       oRsp.Color,
		}
		n.ackChsMu.RUnlock()
	}
}

func (n *Node) commit() {
	var waitingFor uint64
	waitingForMu := sync.Mutex{}

	waitC := make(chan *nodepb.Com)

	cachedCommits := sync.Map{}

	go func() {
		var wF uint64
		for com := range n.possibleComCh {
			waitingForMu.Lock()
			wF = waitingFor
			waitingForMu.Unlock()
			if com.GlobalToken == wF {
				waitC <- com
				continue
			}
			cachedCommits.Store(com.GlobalToken, com)
		}
	}()

	for oRsp := range n.waitingORespCh {
		waitingForMu.Lock()
		waitingFor = oRsp.Gsn
		waitingForMu.Unlock()

		com, ok := cachedCommits.LoadAndDelete(oRsp.Gsn)
		var comMsg *nodepb.Com
		if ok {
			comMsg = com.(*nodepb.Com)
		} else {
			comMsg = <-waitC
		}
		isCoor := n.id == uint32(comMsg.LocalToken)
		if err := n.app.Commit(comMsg.LocalToken, comMsg.Color, comMsg.GlobalToken, isCoor); err != nil {
			logrus.Fatalln("app failed to commit")
		}
	}

}
