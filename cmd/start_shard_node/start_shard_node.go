package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node"
	"github.com/nathanieltornow/PMLog/shared_log"
	"github.com/nathanieltornow/PMLog/shared_log/storage/mem_log"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	IP      = flag.String("IP", ":9000", "")
	logIP   = flag.String("log", ":8000", "")
	peerIPs = flag.String("peers", "", "")
	orderIP = flag.String("order", "", "")
	id      = flag.Int("id", 0, "")
)

func main() {
	flag.Parse()

	log, err := mem_log.NewMemLog()
	if err != nil {
		logrus.Fatalln(err)
	}
	app, err := shared_log.NewSharedLog(log)
	if err != nil {
		logrus.Fatalln(err)
	}
	go func() {
		err := app.Start(*logIP)
		if err != nil {
			logrus.Fatalln(err)
		}
	}()

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
