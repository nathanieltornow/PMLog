package app_node

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/color_service"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	seqpb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
)

func (n *Node) broadcastPrepareMsgs(batchedPrep *nodepb.Prep) {
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
		n.colorServicesMu.RLock()
		cs, ok := n.colorServices[comReq.Color]
		n.colorServicesMu.RUnlock()

		if !ok {
			n.colorServicesMu.Lock()
			cs = color_service.NewColorService(n.id, comReq.Color, n.color, n.app, n.handlePreparation, n.makeOrderRequest)
			n.colorServices[comReq.Color] = cs
			n.colorServicesMu.Unlock()
		}
		cs.PrepareCoordinator(comReq.Content, comReq.FindToken)
	}
}

func (n *Node) makeOrderRequest(oReq *seqpb.OrderRequest) {

}

func (n *Node) handlePreparation(prepMsg *nodepb.Prep) {

}

func (n *Node) handleOrderResponses() {
	for oRsp {

	}
}
