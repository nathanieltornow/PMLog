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

		// put into channel to be broadcasted to other nodes and send orderrequest
		n.prepCh <- prepMsg
		n.orderReqCh <- &seqpb.OrderRequest{Lsn: newToken, Color: comReq.Color, OriginColor: n.color}

		n.prepMan.prepare(prepMsg, comReq.FindToken)
	}
}

func (n *Node) handleOrderResponses() {
	for oRsp := range n.orderRespCh {
		n.prepMan.waitForPrep(oRsp.Lsn)
		if err := n.app.Commit(oRsp.Lsn, oRsp.Color, oRsp.Gsn, uint32(oRsp.Lsn) == n.id); err != nil {
			logrus.Fatalln("failed to commit")
		}
	}
}
