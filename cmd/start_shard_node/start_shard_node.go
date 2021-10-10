package main

import (
	"flag"
	"github.com/nathanieltornow/PMLog/shared_log"
	"github.com/nathanieltornow/PMLog/shared_log/storage/mem_log"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	IP      = flag.String("IP", ":5000", "")
	logIP   = flag.String("log", ":4000", "")
	peerIPs = flag.String("peers", "", "")
	orderIP = flag.String("order", ":9000", "")
	id      = flag.Int("id", 0, "")
)

func main() {
	flag.Parse()

	log, err := mem_log.NewLog()
	if err != nil {
		logrus.Fatalln(err)
	}
	sharedLog, err := shared_log.NewSharedLog(log, uint32(*id), 1010)
	if err != nil {
		logrus.Fatalln(err)
	}

	var peerList []string
	if *peerIPs != "" {
		peerList = strings.Split(*peerIPs, ",")
	}
	err = sharedLog.Start(*logIP, *IP, *orderIP, peerList)
	if err != nil {
		logrus.Fatalln(err)
	}
}
