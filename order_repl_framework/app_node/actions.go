package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/sirupsen/logrus"
)

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

func (n *Node) handleAppCommitRequests() {
	comReqCh := make(chan *frame.CommitRequest, 1024)
	go func() {
		err := n.app.MakeCommitRequests(comReqCh)
		if err != nil {
			return
		}
	}()

	for comReq := range comReqCh {
		n.getLocalNode(comReq.Color).PutComReq(comReq)
	}
}

func (n *Node) handleOrderResponses() {
	for {
		oRsp := n.orderClient.GetNextOrderResponse()
		n.getLocalNode(oRsp.Color).PutOrderResponse(oRsp)
	}
}
