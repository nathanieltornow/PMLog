package app_node

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/sirupsen/logrus"
)

func (n *Node) Prepare(stream nodepb.Node_PrepareServer) error {
	for {
		prepMsg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive prep-msg: %v", err)
		}

		n.getLocalNode(prepMsg.Color).Prepare(prepMsg)
	}
}

func (n *Node) GetAcknowledgements(ackReq *nodepb.AcknowledgeRequest, stream nodepb.Node_GetAcknowledgementsServer) error {
	ackCh := make(chan *nodepb.Acknowledgement, 512)
	n.ackBroadCastMu.Lock()
	n.ackBroadCast[ackReq.ID] = ackCh
	n.ackBroadCastMu.Unlock()

	for ack := range ackCh {
		err := stream.Send(ack)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
	return nil
}

func (n *Node) Register(_ context.Context, req *nodepb.RegisterRequest) (*nodepb.Empty, error) {
	if err := n.connectToPeer(req.IP, false); err != nil {
		return nil, err
	}
	return &nodepb.Empty{}, nil
}
