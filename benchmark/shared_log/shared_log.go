package main

import (
	"flag"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/nathanieltornow/PMLog/benchmark"
	log_client "github.com/nathanieltornow/PMLog/client"
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
	appendLatencies []float64
	readLatencies   []float64
}

type overallResult struct {
	throughputs     []float64
	appendLatencies []float64
	appendMedians   []float64
	readLatencies   []float64
	readMedians     []float64
	append99Perc    []float64
	read99Perc      []float64
	append95Perc    []float64
	read95Perc      []float64
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

	f, err := os.OpenFile(fmt.Sprintf("results_a%v-r%v.csv", config.Appends, config.Reads),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	clients := make([]*log_client.Client, 0)
	for i := 0; i < config.Clients; i++ {
		client, err := log_client.NewClient(uint32(i), config.Endpoints)
		if err != nil {
			logrus.Fatalln(err)
		}
		clients = append(clients, client)
	}

	overall := overallResult{
		throughputs:     make([]float64, 0),
		appendLatencies: make([]float64, 0),
		readLatencies:   make([]float64, 0),
		append99Perc:    make([]float64, 0),
		read99Perc:      make([]float64, 0),
		append95Perc:    make([]float64, 0),
		read95Perc:      make([]float64, 0),
	}

	for t := 0; t < config.Times; t++ {
		resultC = make(chan *benchmarkResult, threads*len(clients))
		for _, client := range clients {
			for i := 0; i < threads; i++ {
				go executeBenchmark(client, config.Runtime, config.Appends, config.Reads)
			}
		}
		appendLatencies := make([]float64, 0)
		readLatencies := make([]float64, 0)
		for i := 0; i < threads*len(clients); i++ {
			res := <-resultC
			appendLatencies = append(appendLatencies, res.appendLatencies...)
			readLatencies = append(readLatencies, res.readLatencies...)
		}

		overallAppendLatency, err := stats.Mean(appendLatencies)
		if err != nil {
			logrus.Fatalln(err)
		}
		overallReadLatency, err := stats.Mean(readLatencies)
		if err != nil {
			logrus.Fatalln(err)
		}

		appendMedian, err := stats.Median(appendLatencies)
		if err != nil {
			logrus.Fatalln(err)
		}
		readMedian, err := stats.Median(readLatencies)
		if err != nil {
			logrus.Fatalln(err)
		}

		append99Percentile, err := stats.Percentile(appendLatencies, 99)
		if err != nil {
			logrus.Fatalln(err)
		}
		read99Percentile, err := stats.Percentile(readLatencies, 99)
		if err != nil {
			logrus.Fatalln(err)
		}

		append95Percentile, err := stats.Percentile(appendLatencies, 95)
		if err != nil {
			logrus.Fatalln(err)
		}
		read95Percentile, err := stats.Percentile(readLatencies, 95)
		if err != nil {
			logrus.Fatalln(err)
		}

		appendThroughput := float64(len(appendLatencies)) / config.Runtime.Seconds()
		readThroughput := float64(len(readLatencies)) / config.Runtime.Seconds()

		overall.throughputs = append(overall.throughputs, float64(appendThroughput+readThroughput))

		overall.appendLatencies = append(overall.appendLatencies, overallAppendLatency)
		overall.appendMedians = append(overall.appendMedians, appendMedian)
		overall.append99Perc = append(overall.append99Perc, append99Percentile)
		overall.append95Perc = append(overall.append95Perc, append95Percentile)

		overall.readLatencies = append(overall.readLatencies, overallReadLatency)
		overall.readMedians = append(overall.readMedians, readMedian)
		overall.read99Perc = append(overall.read99Perc, read99Percentile)
		overall.read95Perc = append(overall.read95Perc, read95Percentile)

		fmt.Printf("-----\nAppend:\nLatency: %v, %v\nThroughput (ops/s): %v\n-----\nRead:\nLatency: %v, %v\nThroughput (ops/s): %v\n",
			overallAppendLatency, appendMedian, appendThroughput, overallReadLatency, readMedian, readThroughput)

		for _, client := range clients {
			client.Trim(0, 0)
		}
	}

	throughput, err := stats.Mean(overall.throughputs)

	appendLatency, err := stats.Mean(overall.appendLatencies)
	appendMedian, err := stats.Mean(overall.appendMedians)
	append99Perc, err := stats.Mean(overall.append99Perc)
	append95Perc, err := stats.Mean(overall.append95Perc)

	readLatency, err := stats.Mean(overall.readLatencies)
	readMedian, err := stats.Mean(overall.readMedians)
	read99Perc, err := stats.Mean(overall.read99Perc)
	read95Perc, err := stats.Mean(overall.read95Perc)
	if err != nil {
		logrus.Fatalln(err)
	}

	if _, err := f.WriteString(fmt.Sprintf("%v, %v, %v, %v, %v, %v, %v, %v, %v\n", throughput, appendLatency, appendMedian, append99Perc, append95Perc, readLatency, readMedian, read99Perc, read95Perc)); err != nil {
		logrus.Fatalln(err)
	}

}

func executeBenchmark(client *log_client.Client, runtime time.Duration, numAppends, numReads int) {

	appendLatencies := make([]float64, 0)
	readLatencies := make([]float64, 0)

	var appends, reads int

	appendHeavy := numAppends > numReads

	defer func() {
		fmt.Println(len(appendLatencies), len(readLatencies))
		resultC <- &benchmarkResult{
			appendLatencies: appendLatencies,
			readLatencies:   readLatencies,
		}
	}()

	if *wait {
		<-time.After(time.Until(time.Now().Truncate(time.Minute).Add(time.Minute)))
	}

	stop := time.After(10 * time.Second)
	loadLoop(client, stop)

	stop = time.After(runtime)

	if appendHeavy {
		ratio := numAppends / numReads
		fmt.Println(ratio)
		for {
			select {
			case <-stop:
				return
			default:
				start := time.Now()
				gsn := client.Append(record, 0)
				appendLatencies = append(appendLatencies, time.Since(start).Seconds())
				appends++
				if (appends+reads)%ratio == 0 {
					start := time.Now()
					_ = client.Read(gsn, 0)
					readLatencies = append(readLatencies, time.Since(start).Seconds())
				}
				reads++
			}
		}
	}
	var gsn uint64
	ratio := numReads / numAppends
	for {
		select {
		case <-stop:
			return
		default:
			if (appends+reads)%ratio == 0 {
				start := time.Now()
				gsn = client.Append(record, 0)
				appendLatencies = append(appendLatencies, time.Since(start).Seconds())
				appends++
			}
			start := time.Now()
			_ = client.Read(gsn, 0)
			readLatencies = append(readLatencies, time.Since(start).Seconds())
			reads++
		}
	}

}

func loadLoop(client *log_client.Client, stop <-chan time.Time) {
	for {
		select {
		case <-stop:
			return
		default:
			gsn := client.Append(record, 0)
			_ = client.Read(gsn, 0)
		}
	}
}
