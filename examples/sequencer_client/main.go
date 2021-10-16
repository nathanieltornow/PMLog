package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/sequencer/client"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
)

var (
	IP    = flag.String("IP", ":8000", "")
	color = flag.Int("color", 0, "")
)

func main() {
	flag.Parse()
	cl, err := client.NewClient(*IP, 123)
	if err != nil {
		logrus.Fatalln(err)
	}
	go func() {
		for i := 0; i < 20; i++ {
			cl.MakeOrderRequest(&sequencerpb.OrderRequest{Color: uint32(*color), OriginColor: 123, NumOfRecords: 5, Tokens: []uint64{132, 424}})
		}
	}()

	for i := 0; i < 20; i++ {
		oRsp := cl.GetNextOrderResponse()
		fmt.Println(oRsp)
	}
}
