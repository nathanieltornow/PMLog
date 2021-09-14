package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	seqpb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
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
		prepMsg := &nodepb.Prep{
			LocalToken: n.getNewLocalToken(),
			Color:      comReq.Color,
			Content:    comReq.Content,
		}
		if err := n.app.Prepare(prepMsg.LocalToken, prepMsg.Color, prepMsg.Content); err != nil {
			return
		}
		// put into channel to be broadcasted to other nodes and send orderrequest
		n.prepCh <- prepMsg
		n.orderCh <- &seqpb.OrderRequest{Lsn: prepMsg.LocalToken, Color: prepMsg.Color, OriginColor: n.color}
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
		if (num + 1) == n.numOfPeers {
			if err := n.app.Commit(ackMsg.LocalToken); err != nil {
				logrus.Fatalln("Failed to commit")
			}
			n.comCh <- &nodepb.Com{LocalToken: ackMsg.LocalToken}
		}
		n.numOfAcksMu.Unlock()
	}
}

func (n *Node) sendOrderRequests() {
	ch := n.orderClient.MakeOrderRequests()
	for oReq := range n.orderCh {
		ch <- oReq
	}
}

func (n *Node) handleOrderResponses() {
	ch := n.orderClient.GetOrderResponses()
	for oRsp := range ch {
		id := uint32(oRsp.Lsn << 32)
		if id == n.id {
			continue
		}
		if !n.app.IsPrepared(oRsp.Lsn) {
			logrus.Fatalln("not prepared")
		}
		n.ackChsMu.RLock()
		n.ackChs[id] <- &nodepb.Ack{LocalToken: oRsp.Lsn}
		n.ackChsMu.RUnlock()
	}
}
