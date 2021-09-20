package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	seqpb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
)

func (n *Node) broadcastPrepareMsgs(batchedPrep *nodepb.BatchedPrep) {
	n.prepStreamsMu.RLock()
	for _, stream := range n.prepStreams {
		if err := stream.Send(batchedPrep); err != nil {
			logrus.Fatalf("failed to send prepmsg: %v", err)
		}
	}
	n.prepStreamsMu.RUnlock()
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

		n.prepB.add(prepMsg)
		// put into channel to be broadcasted to other nodes and send orderrequest
		n.orderClient.MakeOrderRequest(&seqpb.OrderRequest{LocalToken: newToken, Color: comReq.Color, OriginColor: n.color})

		n.prepMan.prepare(prepMsg, comReq.FindToken)
	}
}

func (n *Node) handleOrderResponses() {
	for {
		oRsp := n.orderClient.GetNextOrderResponse()

		n.prepMan.waitForPrep(oRsp.LocalToken)
		if err := n.app.Commit(oRsp.LocalToken, oRsp.Color, oRsp.GlobalToken, uint32(oRsp.LocalToken) == n.id); err != nil {
			logrus.Fatalln("failed to commit")
		}
	}
}
