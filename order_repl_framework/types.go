package order_repl_framework

import (
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
)

type Application interface {
	Prepare(localToken uint64, color uint32, content string) error

	Commit(localToken uint64, color uint32, globalToken uint64) error

	Acknowledge(localToken uint64, color uint32, globalToken uint64) error
}

type CommitRequest struct {
	Color   uint32
	Content string
}

type MakeCommReqFunc func(commReq *CommitRequest) uint64

type OnPrepFunc func(prep *nodepb.Prep)

type OnOrderFunc func(oReq *sequencerpb.OrderRequest)

type OnCommFunc func(ack *nodepb.Acknowledgement)
