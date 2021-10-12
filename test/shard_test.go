package test

import (
	"flag"
	"fmt"
	"github.com/nathanieltornow/PMLog/shared_log/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	shardIPs   = flag.String("IPs", "", "")
	clientList = make([]*client.Client, 0)
	record     = strings.Repeat("p", 4000)
)

func TestMain(m *testing.M) {
	flag.Parse()
	if *shardIPs == "" {
		return
	}
	ipList := strings.Split(*shardIPs, ",")
	for _, ip := range ipList {
		newClient, err := client.NewClient(ip)
		if err != nil {
			logrus.Fatalln("failed to start client", err)
		}
		clientList = append(clientList, newClient)
	}
	ret := m.Run()
	os.Exit(ret)
}

func TestAppendRead(t *testing.T) {
	numOfAppends := 1000
	// pic one appender
	rand.Seed(time.Now().Unix())
	appendClient := clientList[rand.Intn(len(clientList))]

	gsns := make([]uint64, numOfAppends)
	for i := 0; i < numOfAppends; i++ {
		gsn, err := appendClient.Append(0, record)
		require.NoError(t, err)
		gsns[i] = gsn
	}
	for _, cl := range clientList {
		for _, gsn := range gsns {
			rec, err := cl.Read(0, gsn)
			fmt.Println(gsn, rec)
			require.NoError(t, err)
			require.Equal(t, rec, record)
		}
	}
}
