package app_node

import (
	"fmt"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"sync/atomic"
)

func (n *Node) Prepare(stream nodepb.Node_PrepareServer) error {
	for {
		prepMsg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive prep-msg: %v", err)
		}

		n.recentlyPreparedMu.Lock()
		waitC, ok := n.recentlyPrepared[prepMsg.LocalToken]
		if !ok {
			waitC = make(chan bool, 1)
			n.recentlyPrepared[prepMsg.LocalToken] = waitC
		}
		n.recentlyPreparedMu.Unlock()

		if err := n.app.Prepare(prepMsg.LocalToken, prepMsg.Color, prepMsg.Content, 0, waitC); err != nil {
			return err
		}

	}
}

func (n *Node) GetAcks(req *nodepb.AckReq, stream nodepb.Node_GetAcksServer) error {
	ackCh := make(chan *nodepb.Ack, 512)

	n.ackChsMu.Lock()
	n.ackChs[req.NodeID] = ackCh
	n.ackChsMu.Unlock()

	atomic.AddUint32(&n.numOfPeers, 1)

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
		n.possibleComCh <- comMsg
	}
}
