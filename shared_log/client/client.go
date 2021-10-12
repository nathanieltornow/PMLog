package client

import (
	"context"
	pb "github.com/nathanieltornow/PMLog/shared_log/shared_logpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
)

type Client struct {
	mu         sync.Mutex
	tokenCtr   uint32
	pbClient   pb.SharedLogClient
	waitingGsn map[uint32]chan uint64
	appReqCh   chan *pb.AppendRequest
}

func NewClient(IP string) (*Client, error) {
	cl := new(Client)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cl.pbClient = pb.NewSharedLogClient(conn)
	cl.waitingGsn = make(map[uint32]chan uint64)
	cl.appReqCh = make(chan *pb.AppendRequest, 1024)
	stream, err := cl.pbClient.AsyncAppend(context.Background())
	go cl.sendAppends(stream)
	go cl.receiveAppendResponses(stream)
	return cl, nil
}

func (cl *Client) Append(color uint32, record string) (uint64, error) {
	token := atomic.AddUint32(&cl.tokenCtr, 1)
	req := &pb.AppendRequest{Record: record, Color: color, Token: token}
	waitGsn := make(chan uint64, 1)
	cl.mu.Lock()
	cl.waitingGsn[token] = waitGsn
	cl.mu.Unlock()
	cl.appReqCh <- req

	gsn := <-waitGsn
	return gsn, nil
}

func (cl *Client) sendAppends(stream pb.SharedLog_AsyncAppendClient) {
	for appReq := range cl.appReqCh {
		err := stream.Send(appReq)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (cl *Client) receiveAppendResponses(stream pb.SharedLog_AsyncAppendClient) {
	for {
		appRsp, err := stream.Recv()
		if err != nil {
			logrus.Fatalln(err)
		}
		cl.mu.Lock()
		cl.waitingGsn[appRsp.Token] <- appRsp.Gsn
		cl.mu.Unlock()
	}
}

func (cl *Client) Read(color uint32, gsn uint64) (string, error) {
	req := &pb.ReadRequest{Gsn: gsn}
	resp, err := cl.pbClient.Read(context.Background(), req)
	if err != nil {
		return "", err
	}
	return resp.Record, nil
}

func (cl *Client) Trim(color uint32, gsn uint64) error {
	req := &pb.TrimRequest{Gsn: gsn}
	_, err := cl.pbClient.Trim(context.Background(), req)
	if err != nil {
		return err
	}
	return nil
}
