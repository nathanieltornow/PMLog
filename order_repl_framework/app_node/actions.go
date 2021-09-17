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

		newToken := n.getNewLocalToken()

		prepMsg := &nodepb.Prep{
			LocalToken: newToken,
			Color:      comReq.Color,
			Content:    comReq.Content,
		}

		n.prepMan.prepare(prepMsg, comReq.FindToken)

		// put into channel to be broadcasted to other nodes and send orderrequest
		n.prepCh <- prepMsg
		n.orderReqCh <- &seqpb.OrderRequest{Lsn: newToken, Color: comReq.Color, OriginColor: n.color}
	}
}

func (n *Node) receiveAcks(stream nodepb.Node_GetAcksClient) {
	for {
		ackMsg, err := stream.Recv()
		if err != nil {
			logrus.Fatalln("failed to receive ack")
		}
		n.comMan.receiveAck(ackMsg)
	}
}

func (n *Node) handleOrderResponses() {
	for oRsp := range n.orderRespCh {
		n.waitingORespCh <- oRsp

		id := uint32(oRsp.Lsn)
		if id == n.id {
			continue
		}

		n.prepMan.waitForPrep(oRsp.Lsn)

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
	for oRsp := range n.waitingORespCh {
		n.comMan.waitForCom(oRsp.Lsn)

		if uint32(oRsp.Lsn) == n.id {
			n.comCh <- &nodepb.Com{LocalToken: oRsp.Lsn, Color: oRsp.Color, GlobalToken: oRsp.Gsn}
		}
	}
}
