package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/benchmark"
	log_client "github.com/nathanieltornow/PMLog/shared_log/client"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	configPath = flag.String("config", "", "")
	resultC    chan *benchmarkResult
)

type benchmarkResult struct {
	latency    time.Duration
	throughput int
}

func main() {
	flag.Parse()
	if *configPath == "" {
		logrus.Fatalln("no config file")
	}

	config, err := benchmark.GetBenchConfig(*configPath)
	if err != nil {
		logrus.Fatalln(err)
	}

	threads := config.Threads

	resultC = make(chan *benchmarkResult, threads)
	interval := time.Duration(time.Second.Nanoseconds() / int64(config.Ops))

	f, err := os.OpenFile("result.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	for t := threads; t < threads+20; t++ {
		for i := 0; i < t; i++ {
			go benchmarkLog(config.Endpoint, config.Runtime, interval)
		}

		overallThroughput := 0
		latencySum := time.Duration(0)
		for i := 0; i < t; i++ {
			res := <-resultC
			overallThroughput += res.throughput
			latencySum += res.latency
		}
		ovrLatency := time.Duration(latencySum.Nanoseconds() / int64(t))
		throughputPerSecond := float64(overallThroughput) / config.Runtime.Seconds()
		fmt.Printf("-----\nLatency: %v\nThroughput (ops/s): %v\n", ovrLatency, throughputPerSecond)

		if _, err := f.WriteString(fmt.Sprintf("%v, %v\n", throughputPerSecond, ovrLatency.Microseconds())); err != nil {
			logrus.Fatalln(err)
		}
	}
}

func benchmarkLog(IP string, runtime, interval time.Duration) {
	client, err := log_client.NewClient(IP)
	if err != nil {
		logrus.Fatalln(err)
	}

	var latencySum time.Duration
	var throughput int

	defer func() {
		avgLatency := time.Duration(latencySum.Nanoseconds() / int64(throughput))
		resultC <- &benchmarkResult{latency: avgLatency, throughput: throughput}
	}()

	ticker := time.Tick(interval)
	//<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	stop := time.After(runtime)

	for {
		select {
		case <-stop:
			return
		case <-ticker:

			start := time.Now()
			_, err := client.Append(0, "Hello")
			latencySum += time.Since(start)
			throughput++
			if err != nil {
				logrus.Fatalln(err)
			}
		}
	}
}
