package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/sirupsen/logrus"
)

func (n *Node) MakeCommitRequest(comReq *frame.CommitRequest) uint64 {
	return n.getLocalNode(comReq.Color).PutComReq(comReq)
}

func (n *Node) makePrepareMsg(prepMsg *nodepb.Prep) {
	n.prepareCh <- prepMsg
}

func (n *Node) broadcastPrepareMsgs() {
	for prepMsg := range n.prepareCh {
		n.prepStreamsMu.RLock()
		for _, stream := range n.prepStreams {
			if err := stream.Send(prepMsg); err != nil {
				logrus.Fatalf("failed to send prepmsg: %v", err)
			}
		}
		n.prepStreamsMu.RUnlock()
	}

}

func (n *Node) handleOrderResponses() {
	for {
		oRsp := n.orderClient.GetNextOrderResponse()
		n.getLocalNode(oRsp.Color).PutOrderResponse(oRsp)
	}
}

func (n *Node) sendAck(ack *nodepb.Acknowledgement) {
	n.ackBroadCastMu.RLock()
	ackCh, ok := n.ackBroadCast[uint32(ack.LocalToken>>32)]
	if !ok {
		logrus.Fatalln("failed to get ackCh")
	}
	ackCh <- ack
	n.ackBroadCastMu.RUnlock()
}

func (n *Node) handleAcks(stream nodepb.Node_GetAcknowledgementsClient) {
	for {
		ack, err := stream.Recv()
		if err != nil {
			logrus.Fatalln(err)
		}
		n.getLocalNode(ack.Color).PutAcknowledgment(ack)
	}
}
