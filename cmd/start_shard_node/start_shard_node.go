package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/examples/hello_world"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	IP      = flag.String("IP", "", "")
	peerIPs = flag.String("peers", "", "")
	orderIP = flag.String("order", "", "")
	id      = flag.Int("id", 0, "")
)

func main() {
	flag.Parse()

	app := &hello_world.HelloWorldApp{}

	node, err := app_node.NewNode(uint32(*id), 100, app)
	if err != nil {
		logrus.Fatalln(err)
	}
	var peerList []string
	if *peerIPs != "" {
		peerList = strings.Split(*peerIPs, ",")
	}
	err = node.Start(*IP, peerList, *orderIP)
	if err != nil {
		logrus.Fatalln(err)
	}
}
