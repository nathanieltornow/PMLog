package order_repl_framework

type Application interface {
	MakeCommitRequests(chan *CommitRequest) error

	Prepare(localToken uint64, color uint32, content string, findToken uint64, finished chan bool) error

	Commit(localToken uint64, color uint32, globalToken uint64, isCoordinator bool) error
}

type CommitRequest struct {
	Color     uint32
	Content   string
	FindToken uint64
}
