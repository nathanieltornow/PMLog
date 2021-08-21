package sequencer

import (
	"context"
	"flag"
	"fmt"
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"os"
	"strings"
	"testing"
	"time"
)

var (
	endpointList []string
)

func TestMain(m *testing.M) {
	endpoints := flag.String("endpoints", "", "")
	flag.Parse()
	if *endpoints == "" {
		logrus.Fatalln("no endpoints given")
	}
	split := strings.Split(*endpoints, ",")
	endpointList = split
	ret := m.Run()
	os.Exit(ret)
}

func TestSerial(t *testing.T) {
	stream, err := getClientStream(endpointList[0])
	require.NoError(t, err)
	waitC := make(chan bool)
	fmt.Println("hi")
	go func() {
		for {
			oRsp, err := stream.Recv()
			if err != nil {
				close(waitC)
				return
			}
			fmt.Println(oRsp)
		}
	}()
	for i := uint64(0); i < 10; i++ {
		time.Sleep(time.Second)
		err := stream.Send(&pb.OrderRequest{Lsn: i, NumOfRecords: 1, Color: 3, OriginColor: 4})
		require.NoError(t, err)
	}
	time.Sleep(time.Second)
	err = stream.CloseSend()
	require.NoError(t, err)
	<-waitC
}

func getClientStream(IP string) (pb.Sequencer_GetOrderClient, error) {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := pb.NewSequencerClient(conn)
	stream, err := client.GetOrder(context.Background())
	if err != nil {
		return nil, err
	}
	return stream, nil
}
