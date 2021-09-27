package order_repl_framework

import (
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
)

type Application interface {
	MakeCommitRequests(chan *CommitRequest) error

	Prepare(localToken uint64, color uint32, content string, findToken uint64) error

	Commit(localToken uint64, color uint32, globalToken uint64, isCoordinator bool) error
}

type CommitRequest struct {
	Color     uint32
	Content   string
	FindToken uint64
}

type OnPrepBatch func(prep *nodepb.Prep)

type OnOrderBatch func(prep *sequencerpb.OrderRequest)
