package order_repl_framework

type Application interface {
	MakeCommitRequests(chan *CommitRequest) error

	Prepare(localToken uint64, color uint32, content string) error

	IsPrepared(localToken uint64) bool

	Commit(localToken uint64, color uint32, uint642 uint64) error
}

type CommitRequest struct {
	Color   uint32
	Content string
}
