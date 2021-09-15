package redis

import frame "github.com/nathanieltornow/PMLog/order_repl_framework"

func (s *Server) MakeCommitRequests(ch chan *frame.CommitRequest) error {

	return nil
}

func (s *Server) Prepare(localToken uint64, color uint32, content string, findToken uint64, finished chan bool) error {
	return nil
}

func (s *Server) Commit(localToken uint64, color uint32, globalToken uint64, isCoordinator bool) error {
	return nil
}
