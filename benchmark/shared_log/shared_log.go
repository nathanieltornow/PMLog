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
	avgLatency time.Duration
	numAppends int
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

	numEndpoints := len(config.Endpoints)
	interval := time.Duration(time.Second.Nanoseconds() / int64(config.Ops))

	f, err := os.OpenFile("result.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	for t := threads; t < threads+20; t++ {
		resultC = make(chan *benchmarkResult, t*numEndpoints)
		for _, endpoint := range config.Endpoints {
			for i := 0; i < t; i++ {
				go benchmarkLog(endpoint, config.Runtime, interval)
			}
		}
		overallAppends := 0
		weightedLatencySum := time.Duration(0)
		for i := 0; i < t*numEndpoints; i++ {
			res := <-resultC
			overallAppends += res.numAppends
			weightedLatencySum += time.Duration(res.avgLatency.Nanoseconds() * int64(res.numAppends))
		}
		ovrLatency := time.Duration(weightedLatencySum.Nanoseconds() / int64(overallAppends))
		throughputPerSecond := float64(overallAppends) / config.Runtime.Seconds()
		fmt.Printf("-----\nLatency: %v\nThroughput (ops/s): %v\n", ovrLatency, throughputPerSecond)
		if _, err := f.WriteString(fmt.Sprintf("%v, %v\n", throughputPerSecond, ovrLatency)); err != nil {
			logrus.Fatalln(err)
		}
	}

}

func benchmarkLog(IP string, runtime, interval time.Duration) {
	client, err := log_client.NewClient(IP)
	if err != nil {
		logrus.Fatalln(err)
	}

	overallLatency := time.Duration(0)
	var appends int

	defer func() {
		avgLatency := time.Duration(overallLatency.Nanoseconds() / int64(appends))
		resultC <- &benchmarkResult{avgLatency: avgLatency, numAppends: appends}
	}()

	ticker := time.Tick(interval)
	<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	stop := time.After(runtime)
	for {
		select {
		case <-stop:
			return
		case <-ticker:
			start := time.Now()
			_, err := client.Append(0, "Hello")
			overallLatency += time.Since(start)
			if err != nil {
				logrus.Fatalln(err)
			}
			appends++
		}
	}
}

func getAvg(latencies []time.Duration) int {
	sum := time.Duration(0)
	for _, lat := range latencies {
		sum += lat
	}
	avg := sum.Microseconds() / int64(len(latencies))
	return int(avg)
}
