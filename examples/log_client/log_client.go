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

	for i := 0; i < 10; i++ {
		time.Sleep(3 * time.Second)
		resp, err := logClient.Append(0, "Hello")
		if err != nil {
			logrus.Fatalln(err)
		}
		fmt.Println(resp)
	}
}
