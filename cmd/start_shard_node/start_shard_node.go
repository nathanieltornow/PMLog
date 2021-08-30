package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/shard"
	"github.com/nathanieltornow/PMLog/storage/mem_log"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	IP      = flag.String("IP", "", "")
	peerIPs = flag.String("peers", "", "")
	id      = flag.Int("id", 0, "")
)

func main() {
	flag.Parse()
	primLog, err := mem_log.NewMemLog()
	if err != nil {
		logrus.Fatalln(err)
	}
	secLog, err := mem_log.NewMemLog()
	if err != nil {
		logrus.Fatalln(err)
	}
	node, err := shard.NewNode(uint32(*id), 0, primLog, secLog)
	if err != nil {
		logrus.Fatalln(err)
	}

	var peerList []string
	if *peerIPs != "" {
		peerList = strings.Split(*peerIPs, ",")
	}
	err = node.Start(*IP, peerList)
	if err != nil {
		logrus.Fatalln(err)
	}
}
