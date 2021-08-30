package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/shard_client"
	"github.com/sirupsen/logrus"
	"time"
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
	iter := 100
	dur := time.Duration(0)
	for i := 0; i < iter; i++ {
		start := time.Now()
		_, err := client.Append(*append, 0)
		if err != nil {
			logrus.Fatalln(err)
		}
		dur += time.Since(start)
	}
	fmt.Println(dur.Microseconds() / int64(iter))
}
