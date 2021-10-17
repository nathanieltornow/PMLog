package main

import (
	"flag"
	"fmt"
	"github.com/montanaflynn/stats"
	"github.com/nathanieltornow/PMLog/benchmark"
	seq_client "github.com/nathanieltornow/PMLog/sequencer/client"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	configPath  = flag.String("config", "", "")
	resultC     chan *benchmarkResult
	color       = flag.Int("color", 0, "")
	originColor = flag.Int("origincolor", 10, "")
)

type benchmarkResult struct {
	operations int
	latencySum time.Duration
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

	f, err := os.OpenFile("results_sequencer.csv",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalln(err)
	}
	defer f.Close()

	throughputResults := make([]float64, 0)
	latencyResults := make([]float64, 0)

	for t := 0; t < config.Times; t++ {
		resultC = make(chan *benchmarkResult, config.Clients)
		clients := make([]*seq_client.Client, 0)
		for i := 0; i < config.Clients; i++ {
			cl, err := seq_client.NewClient(config.Endpoints[0], uint32(*originColor))
			if err != nil {
				logrus.Fatalln(err)
			}
			clients = append(clients, cl)
		}
		for _, cl := range clients {
			go executeBenchmark(cl, uint32(*color), uint32(*originColor), config.Runtime)
		}

		overallOperations := 0
		overallLatency := time.Duration(0)

		for i := 0; i < config.Clients; i++ {
			res := <-resultC
			overallOperations += res.operations
			overallLatency += res.latencySum
		}

		throughput := float64(overallOperations) / config.Runtime.Seconds()
		latency := time.Duration(overallLatency.Nanoseconds() / int64(overallOperations))

		throughputResults = append(throughputResults, throughput)
		latencyResults = append(latencyResults, latency.Seconds())

		fmt.Printf("-----Latency: %v, \nThroughput (ops/s): %v\n", latency, throughput)
	}

	overallThroughput, err := stats.Mean(throughputResults)
	overallLatency, err := stats.Mean(latencyResults)
	if err != nil {
		logrus.Fatalln(err)
	}

	if _, err := f.WriteString(fmt.Sprintf("%v, %v\n", overallThroughput, overallLatency)); err != nil {
		logrus.Fatalln(err)
	}

}

func executeBenchmark(client *seq_client.Client, color, originColor uint32, duration time.Duration) {
	stop := time.After(duration)
	waitC := make(chan bool, 1)

	latencySum := time.Duration(0)
	operations := 0

	defer func() {
		resultC <- &benchmarkResult{operations: operations, latencySum: latencySum}
	}()

	go func() {
		for {
			select {
			case <-stop:
				waitC <- true
				return
			default:
				client.MakeOrderRequest(&sequencerpb.OrderRequest{OriginColor: originColor, Color: color, NumOfRecords: 1, Tokens: []uint64{uint64(time.Now().UnixNano())}})
			}
		}
	}()
	for {
		select {
		case <-waitC:
			return
		default:
			startTime := client.GetNextOrderResponse().Tokens[0]
			//fmt.Println(startTime)
			latencySum += time.Duration(time.Now().UnixNano() - int64(startTime))
			operations++
		}
	}
}
