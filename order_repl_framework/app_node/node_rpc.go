package app_node

import (
	"fmt"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
)

func (n *Node) Prepare(stream nodepb.Node_PrepareServer) error {
	for {
		batchedPrepMsg, err := stream.Recv()
		if err != nil {
			return fmt.Errorf("failed to receive prep-msg: %v", err)
		}

		for _, prepMsg := range batchedPrepMsg.Preps {
			n.prepMan.prepare(prepMsg, 0)
		}
	}
}
