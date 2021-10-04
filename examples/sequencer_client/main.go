package main

//import (
//	"flag"
//	seqclient "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/client"
//	pb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
//	"github.com/sirupsen/logrus"
//	"time"
//)
//
//var (
//	IP    = flag.String("IP", ":8000", "")
//	color = flag.Int("color", 0, "")
//)
//
//func main() {
//	flag.Parse()
//	waitC := make(chan bool)
//	client, err := seqclient.NewClient(*IP)
//	if err != nil {
//		logrus.Fatalln(err)
//	}
//	go func() {
//		for {
//			oRsp := client.GetNextOrderResponse()
//			logrus.Infof("Got OrderResponse: %v\n", oRsp)
//		}
//
//	}()
//	for i := 0; i < 10; i++ {
//		time.Sleep(10 *time.Second)
//		oReq := &pb.OrderRequest{Lsn: uint64(i), Color: uint32(*color), NumOfRecords: 1}
//		logrus.Infof("Sending OrderRequest %v\n", oReq)
//		client.MakeOrderRequest(oReq)
//	}
//	<-waitC
//	time.Sleep(2 * time.Second)
//	if err != nil {
//		logrus.Fatalln(err)
//	}
//
//}
