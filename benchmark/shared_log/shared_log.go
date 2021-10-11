package main

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/benchmark"
	log_client "github.com/nathanieltornow/PMLog/shared_log/client"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var (
	configPath  = flag.String("config", "", "")
	resultC     chan *benchmarkResult
	record      = strings.Repeat("r", 4000)
	threadsFlag = flag.Int("threads", 0, "")
)

type benchmarkResult struct {
	avgAppendLatency time.Duration
	appends          int
	avgReadLatency   time.Duration
	reads            int
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
	if *threadsFlag != 0 {
		threads = *threadsFlag
	}

	numEndpoints := len(config.Endpoints)
	appendInterval := time.Duration(time.Second.Nanoseconds() / int64(config.Appends))
	readInterval := time.Duration(time.Second.Nanoseconds() / int64(config.Reads))

	f, err := os.OpenFile("result.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	for t := 0; t < config.Times; t++ {
		resultC = make(chan *benchmarkResult, threads*numEndpoints)
		for _, endpoint := range config.Endpoints {
			for i := 0; i < threads; i++ {
				go executeBenchmark(endpoint, config.Runtime, appendInterval, readInterval)
			}
		}

		overallAppends := 0
		weightedAppendLatencySum := time.Duration(0)
		overallReads := 0
		weightedReadLatencySum := time.Duration(0)
		for i := 0; i < threads*numEndpoints; i++ {
			res := <-resultC
			overallAppends += res.appends
			weightedAppendLatencySum += time.Duration(res.avgAppendLatency.Nanoseconds() * int64(res.appends))
			overallReads += res.reads
			weightedReadLatencySum += time.Duration(res.avgReadLatency.Nanoseconds() * int64(res.reads))
		}

		overallAppendLatency := time.Duration(0)
		appendThroughput := float64(0)
		if overallAppends != 0 {
			overallAppendLatency = time.Duration(weightedAppendLatencySum.Nanoseconds() / int64(overallAppends))
			appendThroughput = float64(overallAppends) / config.Runtime.Seconds()
		}

		overallReadLatency := time.Duration(0)
		readThroughput := float64(0)
		if overallReads != 0 {
			overallReadLatency = time.Duration(weightedReadLatencySum.Nanoseconds() / int64(overallReads))
			readThroughput = float64(overallReads) / config.Runtime.Seconds()
		}

		fmt.Printf(
			"-----\nAppend:\nLatency: %v\nThroughput (ops/s): %v\n-----\nRead:\nLatency: %v\nThroughput (ops/s): %v\n",
			overallAppendLatency, appendThroughput, overallReadLatency, readThroughput)

		if _, err := f.WriteString(fmt.Sprintf("%v, %v, %v, %v \n", overallAppendLatency.Microseconds(),
			appendThroughput, overallReadLatency.Microseconds(), readThroughput)); err != nil {
			logrus.Fatalln(err)
		}
	}

}

func executeBenchmark(IP string, runtime, appendInterval, readInterval time.Duration) {
	client, err := log_client.NewClient(IP)
	if err != nil {
		logrus.Fatalln(err)
	}

	var curGsn uint64

	overallAppendLatency := time.Duration(0)
	overallReadLatency := time.Duration(0)
	var appends int
	var reads int

	defer func() {
		avgAppendLatency := time.Duration(0)
		avgReadLatency := time.Duration(0)
		if reads > 0 {
			avgReadLatency = time.Duration(overallReadLatency.Nanoseconds() / int64(reads))
		}
		if appends > 0 {
			avgAppendLatency = time.Duration(overallAppendLatency.Nanoseconds() / int64(appends))
		}

		resultC <- &benchmarkResult{
			avgAppendLatency: avgAppendLatency,
			appends:          appends,
			avgReadLatency:   avgReadLatency,
			reads:            reads,
		}
		err := client.Trim(0, 0)
		if err != nil {
			logrus.Fatalln(err)
		}
	}()

	appendTicker := time.Tick(appendInterval)
	readTicker := time.Tick(readInterval)
	<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	stop := time.After(runtime)
	for {
		select {
		case <-stop:
			return
		case <-appendTicker:
			start := time.Now()
			gsn, err := client.Append(0, record)
			overallAppendLatency += time.Since(start)
			if err != nil {
				logrus.Fatalln(err)
			}
			appends++
			curGsn = gsn
		case <-readTicker:
			if curGsn == 0 {
				continue
			}
			start := time.Now()
			_, err := client.Read(0, curGsn)
			overallReadLatency += time.Since(start)
			if err != nil {
				logrus.Fatalln(err)
			}
			reads++
		}
	}
}
