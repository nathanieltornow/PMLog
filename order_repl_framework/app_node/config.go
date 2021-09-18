package app_node

import "time"

type NodeOption func(node *Node)

func WithPrepBatchInterval(interval time.Duration) NodeOption {
	return func(node *Node) {
		node.batchInterval = interval
	}
}

func WithMaxPrepMsgSize(size int) NodeOption {
	return func(node *Node) {
		node.maxMsgSize = size
	}
}
