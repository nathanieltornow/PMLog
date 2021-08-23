package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/benchmark"
	"github.com/nathanieltornow/PMLog/sequencer"
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	configFile = flag.String("config", "", "")
	colorFlag  = flag.Int("color", 0, "")

	startTime = time.Now() //.Truncate(time.Minute).Add(time.Minute)
	resultC   chan *benchmarkResult
)

func main() {
	flag.Parse()
	config, err := benchmark.GetBenchConfig(*configFile)
	if err != nil {
		logrus.Fatalln(err)
	}
	resultC = make(chan *benchmarkResult, config.Threads)
	interval := time.Duration(time.Second.Nanoseconds() / int64(config.Ops))
	for i := 0; i < config.Threads; i++ {
		go benchmarkSequencer(config.Endpoint, uint32(*colorFlag), uint32(i), config.Runtime, interval)
	}
	overallThroughput := 0
	latencySum := time.Duration(0)
	for i := 0; i < config.Threads; i++ {
		res := <-resultC
		overallThroughput += res.throughput
		latencySum += res.latency
	}
	ovrLatency := time.Duration(latencySum.Nanoseconds() / int64(config.Threads))
	throughputPerSecond := float64(overallThroughput) / config.Runtime.Seconds()
	fmt.Printf("Latency: %v\nThroughput: %v\n", ovrLatency, throughputPerSecond)
}

type benchmarkResult struct {
	latency    time.Duration
	throughput int
}

func benchmarkSequencer(IP string, color, originColor uint32, runtime, interval time.Duration) {
	client, err := sequencer.NewClient(IP)
	if err != nil {
		logrus.Fatalln("failed to start client for sequencer")
	}
	ticker := time.Tick(interval)

	oReqC := client.MakeOrderRequests()
	oRspC := client.GetOrderResponses()

	inTimes := make(map[uint64]time.Time)
	outTimes := make(map[uint64]time.Time)

	go func() {
		for {
			rsp, ok := <-oRspC
			now := time.Now()
			if !ok {
				return
			}
			if rsp.OriginColor == originColor {
				inTimes[rsp.Lsn] = now
			}
		}
	}()

	i := 0

	defer func() {
		time.Sleep(time.Second)
		_ = client.Stop()
		lat := getOverallLatency(inTimes, outTimes)
		resultC <- &benchmarkResult{latency: lat, throughput: i}
	}()

	<-time.After(time.Until(startTime))

	stop := time.After(runtime)

	for {
		select {
		case <-stop:
			return
		case _, ok := <-ticker:
			if !ok {
				return
			}
			oReqC <- &pb.OrderRequest{
				Lsn:          uint64(i),
				NumOfRecords: 1,
				Color:        color,
				OriginColor:  originColor,
			}
			now := time.Now()
			outTimes[uint64(i)] = now
			i++
		}
	}
}

func getOverallLatency(in, out map[uint64]time.Time) time.Duration {
	numOfRes := 0
	var latencySum time.Duration
	for outLsn, outTime := range out {
		inTime, ok := in[outLsn]
		if !ok {
			continue
		}
		numOfRes++
		latencySum += inTime.Sub(outTime)
	}
	return time.Duration(latencySum.Nanoseconds() / int64(numOfRes))
}