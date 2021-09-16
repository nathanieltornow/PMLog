package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	seqpb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"sync/atomic"
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

		single := atomic.LoadUint32(&n.numOfPeers) == 0

		newToken := n.getNewLocalToken()

		n.recentlyPreparedMu.Lock()
		waitC, ok := n.recentlyPrepared[newToken]
		if !ok {
			waitC = make(chan bool, 1)
			n.recentlyPrepared[newToken] = waitC
		}
		n.recentlyPreparedMu.Unlock()

		if err := n.app.Prepare(newToken, comReq.Color, comReq.Content, comReq.FindToken, waitC); err != nil {
			return
		}

		if single {
			<-waitC
			if err := n.app.Commit(newToken, comReq.Color, newToken, true); err != nil {
				logrus.Fatalln("failed to commit")
			}
			continue
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

		if (num + 1) == atomic.LoadUint32(&n.numOfPeers) {
			comMsg := &nodepb.Com{
				LocalToken:  ackMsg.LocalToken,
				Color:       ackMsg.Color,
				GlobalToken: ackMsg.GlobalToken,
			}

			n.recentlyPreparedMu.Lock()
			waitC, ok := n.recentlyPrepared[ackMsg.LocalToken]
			if !ok {
				logrus.Fatalln("Failed to find waitC")
			}
			n.recentlyPreparedMu.Unlock()
			<-waitC

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

		n.recentlyPreparedMu.Lock()
		waitC, ok := n.recentlyPrepared[oRsp.Lsn]
		if !ok {
			waitC = make(chan bool, 1)
			n.recentlyPrepared[oRsp.Lsn] = waitC
		}
		n.recentlyPreparedMu.Unlock()
		<-waitC

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

	cachedCommits := make(map[uint64]*nodepb.Com)

	for oRsp := range n.waitingORespCh {

		com, ok := cachedCommits[oRsp.Gsn]
		if !ok {
			for comMsg := range n.possibleComCh {
				if comMsg.GlobalToken == oRsp.Gsn {
					com = comMsg
					break
				}
				cachedCommits[comMsg.GlobalToken] = comMsg
			}
		}

		isCoor := n.id == uint32(com.LocalToken)
		if err := n.app.Commit(com.LocalToken, com.Color, com.GlobalToken, isCoor); err != nil {
			logrus.Fatalln("app failed to commit")
		}
	}

}
