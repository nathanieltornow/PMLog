package app_node

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
)

func (n *Node) Prepare(stream nodepb.Node_PrepareServer) error {
	for {
		prepMsg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive prep-msg: %v", err)
		}
		if err := n.app.Prepare(prepMsg.LocalToken, prepMsg.Color, prepMsg.Content); err != nil {
			return err
		}
	}
}

func (n *Node) GetAcks(req *nodepb.AckReq, stream nodepb.Node_GetAcksServer) error {
	ackCh := make(chan *nodepb.Ack, 512)
	n.ackChsMu.Lock()
	n.ackChs[req.NodeID] = ackCh
	n.ackChsMu.Unlock()
	for ack := range ackCh {
		if err := stream.Send(ack); err != nil {
			return err
		}
	}
	return nil
}

func (n *Node) Commit(stream nodepb.Node_CommitServer) error {
	for {
		comMsg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive com-msg: %v", err)
		}
		if err := n.app.Commit(comMsg.LocalToken); err != nil {
			return err
		}
	}
}

func (n *Node) Register(_ context.Context, req *nodepb.RegisterRequest) (*nodepb.Empty, error) {
	if err := n.connectToPeer(req.IP, false); err != nil {
		return nil, err
	}
	return nil, nil
}
