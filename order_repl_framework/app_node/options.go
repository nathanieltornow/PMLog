package app_node

import "time"

type NodeOption func(node *Node)

func WithBatchingInterval(interval time.Duration) NodeOption {
	return func(node *Node) {
		node.interval = interval
	}
}
