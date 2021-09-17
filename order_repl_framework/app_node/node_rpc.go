package app_node

import (
	"fmt"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
)

func (n *Node) Prepare(stream nodepb.Node_PrepareServer) error {
	for {
		prepMsg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive prep-msg: %v", err)
		}

		n.prepMan.prepare(prepMsg, 0)
	}
}

func (n *Node) GetAcks(req *nodepb.AckReq, stream nodepb.Node_GetAcksServer) error {
	ackCh := make(chan *nodepb.Ack, 512)

	n.ackChsMu.Lock()
	n.ackChs[req.NodeID] = ackCh
	n.ackChsMu.Unlock()

	n.comMan.addPeer()

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
		fmt.Println("recv com", comMsg)
		if err != nil {
			return fmt.Errorf("failed to receive com-msg: %v", err)
		}
		n.comMan.receiveCom(comMsg)
	}
}
