package hello_world

import (
	"fmt"
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"time"
)

type HelloWorldApp struct{}

func (h *HelloWorldApp) MakeCommitRequests(ch chan *frame.CommitRequest) error {
	for {
		time.Sleep(10 * time.Second)
		ch <- &frame.CommitRequest{
			Color:   0,
			Content: "Hello World",
		}
	}
}

func (h *HelloWorldApp) Prepare(localToken uint64, color uint32, content string) error {
	fmt.Println("preparing", localToken, color, content)
	return nil
}

func (h *HelloWorldApp) IsPrepared(_ uint64) bool {
	return true
}

func (h *HelloWorldApp) Commit(localToken uint64, color uint32, globalToken uint64) error {
	fmt.Println("Committing", localToken, color, globalToken)
	return nil
}
