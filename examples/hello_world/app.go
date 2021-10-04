package main

import (
	"flag"
	"fmt"
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	IP      = flag.String("IP", ":9000", "")
	logIP   = flag.String("log", ":8000", "")
	peerIPs = flag.String("peers", "", "")
	orderIP = flag.String("order", "", "")
	id      = flag.Int("id", 0, "")
	active  = flag.Bool("a", false, "")
)

func main() {
	flag.Parse()

	hw := NewHW()
	if *active {
		time.Sleep(2 * time.Second)

		hw.doStuff()
	}
	time.Sleep(100 * time.Second)
}

type HelloWorldApp struct {
	node *app_node.Node
}

func NewHW() *HelloWorldApp {
	hw := new(HelloWorldApp)
	node, err := app_node.NewNode(uint32(*id), 11)
	if err != nil {
		logrus.Fatalln(err)
	}
	hw.node = node
	node.RegisterApp(hw)

	go node.Start(*IP, strings.Split(*peerIPs, ","), *orderIP)
	return hw
}

func (h *HelloWorldApp) doStuff() {
	for i := 0; i < 40; i++ {
		h.node.MakeCommitRequest(&frame.CommitRequest{Color: 0, Content: "hallo"})
		time.Sleep(5 * time.Millisecond)
	}
	time.Sleep(10 * time.Second)
}

func (h *HelloWorldApp) Prepare(localToken uint64, color uint32, content string) error {
	fmt.Println("preparing", localToken, color, content)
	return nil
}

func (h *HelloWorldApp) Commit(localToken uint64, color uint32, globalToken uint64) error {
	fmt.Println("Committing", localToken, color, globalToken)
	return nil
}

func (h *HelloWorldApp) Acknowledge(localToken uint64, color uint32, globalToken uint64) error {
	fmt.Println("acknowledgin", localToken, color, globalToken)
	return nil
}
