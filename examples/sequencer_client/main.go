package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/sequencer/client"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"time"
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
	iter := 20000
	start := time.Now()
	for i := 0; i < iter; i++ {
		cl.MakeOrderRequest(&sequencerpb.OrderRequest{Color: uint32(*color), OriginColor: 123, NumOfRecords: 12, Tokens: []uint64{12, 1}})
		cl.GetNextOrderResponse()
	}
	fmt.Println(time.Duration(time.Since(start).Nanoseconds() / int64(iter)))
}
