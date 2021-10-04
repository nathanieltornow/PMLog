package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/shared_log"
	"github.com/nathanieltornow/PMLog/shared_log/storage/dummy_log"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	IP      = flag.String("IP", "", "")
	logIP   = flag.String("log", "", "")
	peerIPs = flag.String("peers", "", "")
	orderIP = flag.String("order", "", "")
	id      = flag.Int("id", 0, "")
)

func main() {
	flag.Parse()

	log, err := dummy_log.NewDummyLog()
	if err != nil {
		logrus.Fatalln(err)
	}
	sharedLog, err := shared_log.NewSharedLog(log, uint32(*id), 1010)
	if err != nil {
		logrus.Fatalln(err)
	}
	var peerList []string
	endPeerList := make([]string, 0)
	if *peerIPs != "" {
		peerList = strings.Split(*peerIPs, ",")
		for _, peer := range peerList {
			endPeerList = append(endPeerList, peer+":5000")
		}
	}
	fmt.Println(peerList)
	err = sharedLog.Start((*logIP)+":4000", *IP, (*orderIP)+":7000", peerList)
	if err != nil {
		logrus.Fatalln(err)
	}
}
