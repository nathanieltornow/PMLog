package shared_log

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
)

func (sl *SharedLog) MakeCommitRequests(ch chan *frame.CommitRequest) error {

	return nil
}

func (sl *SharedLog) Prepare(localToken uint64, color uint32, content string) error {
	return nil
}

func (sl *SharedLog) IsPrepared(localToken uint64) bool {
	return true
}

func (sl *SharedLog) Commit(localToken uint64, color uint32, uint642 uint64) error {
	return nil
}
