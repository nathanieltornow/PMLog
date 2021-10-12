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
	wait        = flag.Bool("wait", false, "")
)

type benchmarkResult struct {
	overallAppendLatency time.Duration
	appends              int
	overallReadLatency   time.Duration
	reads                int
}

type overallResult struct {
	throughput    int
	appendLatency time.Duration
	readLatency   time.Duration
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

	appendInterval := time.Duration(time.Second.Nanoseconds() / int64(config.Appends))
	readInterval := time.Duration(time.Second.Nanoseconds() / int64(config.Reads))

	f, err := os.OpenFile(fmt.Sprintf("results_a%v-r%v.csv", config.Appends, config.Reads),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	clients := make([]*log_client.Client, 0)
	for _, endpoint := range config.Endpoints {
		for i := 0; i < config.Clients; i++ {
			client, err := log_client.NewClient(endpoint)
			if err != nil {
				logrus.Fatalln(err)
			}
			clients = append(clients, client)
		}
	}

	result := &overallResult{}

	for t := 0; t < config.Times; t++ {
		resultC = make(chan *benchmarkResult, threads*len(clients))
		for _, client := range clients {
			for i := 0; i < threads; i++ {
				go executeBenchmark(client, config.Runtime, appendInterval, readInterval)
			}
		}
		overallAppends := 0
		overallAppendLatencySum := time.Duration(0)
		overallReads := 0
		overallReadLatencySum := time.Duration(0)
		for i := 0; i < threads*len(clients); i++ {
			res := <-resultC
			overallAppends += res.appends
			overallAppendLatencySum += res.overallAppendLatency
			overallReads += res.reads
			overallReadLatencySum += res.overallReadLatency
		}

		overallAppendLatency := time.Duration(0)
		appendThroughput := float64(0)
		if overallAppends != 0 {
			overallAppendLatency = time.Duration(overallAppendLatencySum.Nanoseconds() / int64(overallAppends))
			appendThroughput = float64(overallAppends) / config.Runtime.Seconds()
		}

		overallReadLatency := time.Duration(0)
		readThroughput := float64(0)
		if overallReads != 0 {
			overallReadLatency = time.Duration(overallReadLatencySum.Nanoseconds() / int64(overallReads))
			readThroughput = float64(overallReads) / config.Runtime.Seconds()
		}

		result.throughput += int(readThroughput + appendThroughput)
		result.appendLatency += overallAppendLatency
		result.readLatency += overallReadLatency

		fmt.Printf(
			"-----\nAppend:\nLatency: %v\nThroughput (ops/s): %v\n-----\nRead:\nLatency: %v\nThroughput (ops/s): %v\n",
			overallAppendLatency, appendThroughput, overallReadLatency, readThroughput)
	}

	throughput := result.throughput / config.Times
	appendLatency := time.Duration(result.appendLatency.Nanoseconds() / int64(config.Times))
	readLatency := time.Duration(result.readLatency.Nanoseconds() / int64(config.Times))

	if _, err := f.WriteString(fmt.Sprintf("%v, %v, %v\n", throughput,
		appendLatency.Microseconds(), readLatency.Microseconds())); err != nil {
		logrus.Fatalln(err)
	}

}

func executeBenchmark(client *log_client.Client, runtime, appendInterval, readInterval time.Duration) {
	var curGsn uint64
	overallAppendLatency := time.Duration(0)
	overallReadLatency := time.Duration(0)
	var appends int
	var reads int

	defer func() {

		resultC <- &benchmarkResult{
			overallAppendLatency: overallAppendLatency,
			appends:              appends,
			overallReadLatency:   overallReadLatency,
			reads:                reads,
		}
		err := client.Trim(0, 0)
		if err != nil {
			logrus.Fatalln(err)
		}
	}()
	appendTicker := time.Tick(appendInterval)
	readTicker := time.Tick(readInterval)

	if *wait {
		<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	}

	stop := time.After(5 * time.Second)
	loadLoop(client, stop, appendTicker, readTicker)

	stop = time.After(runtime)
benchLoop:
	for {
		select {
		case <-stop:
			break benchLoop
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
	stop = time.After(5 * time.Second)
	loadLoop(client, stop, appendTicker, readTicker)
}

func loadLoop(client *log_client.Client, stop <-chan time.Time, appendTicker, readTicker <-chan time.Time) {
	curGsn := uint64(0)
	for {
		select {
		case <-stop:
			return
		case <-appendTicker:
			gsn, err := client.Append(0, record)
			if err != nil {
				logrus.Fatalln(err)
			}
			curGsn = gsn
		case <-readTicker:
			if curGsn == 0 {
				continue
			}
			_, err := client.Read(0, curGsn)
			if err != nil {
				logrus.Fatalln(err)
			}
		}
	}
}
