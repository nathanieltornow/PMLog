package main

import (
	"flag"
	seqclient "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/client"
	pb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	IP    = flag.String("IP", ":8000", "")
	color = flag.Int("color", 0, "")
)

func main() {
	flag.Parse()
	client, err := seqclient.NewClient(*IP)
	if err != nil {
		logrus.Fatalln(err)
	}

	go func() {
		oRspC := client.GetOrderResponses()
		for oRsp := range oRspC {
			logrus.Infof("Got OrderResponse: %v\n", oRsp)
		}
	}()

	oReqC := client.MakeOrderRequests()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second)
		oReq := &pb.OrderRequest{Lsn: uint64(i), Color: uint32(*color)}
		logrus.Infof("Sending OrderRequest %v\n", oReq)
		oReqC <- oReq
	}
	time.Sleep(2 * time.Second)
	err = client.Stop()
	if err != nil {
		logrus.Fatalln(err)
	}

}
