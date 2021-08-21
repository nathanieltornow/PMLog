package main

import (
	"flag"
	"fmt"
	"time"
)

var (
	parentIP = flag.String("parIP", "", "")
)

func stuff(m map[uint64]chan bool) {
	for i := 0; i < 1000; i++ {
		m[1] <- true
	}
}

func ok(ch chan bool) {
	for {
		time.Sleep(time.Second)
		ch <- true
	}
}

func main() {
	ticker := time.Tick(2 * time.Second)
	ch := make(chan bool, 10)
	go ok(ch)
	for {
		select {
		case <-ticker:
			fmt.Println("hi")
		case <-ch:
			fmt.Println("hoo")
		}

	}
}
