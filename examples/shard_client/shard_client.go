package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/shard_client"
	"github.com/sirupsen/logrus"
)

var (
	IP     = flag.String("IP", "", "")
	append = flag.String("append", "", "")
)

func main() {
	flag.Parse()
	client, err := shard_client.NewClient(*IP)
	if err != nil {
		logrus.Fatalln(err)
	}
	gsn, err := client.Append(*append, 0)
	fmt.Println(gsn)
}
