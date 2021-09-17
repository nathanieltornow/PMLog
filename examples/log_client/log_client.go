package main

import (
	"flag"
	"fmt"
	log_client "github.com/nathanieltornow/PMLog/shared_log/client"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	IP = flag.String("IP", ":8000", "")
)

func main() {
	flag.Parse()

	logClient, err := log_client.NewClient(*IP)
	if err != nil {
		logrus.Fatalln(err)
	}

	for i := 0; i < 20; i++ {
		now := time.Now()
		resp, err := logClient.Append(0, "Hello")
		end := time.Since(now)
		if err != nil {
			logrus.Fatalln(err)
		}
		fmt.Println(resp, end)
	}
}
