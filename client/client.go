package client

import (
	"context"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"sync/atomic"
)

type Client struct {
	mu sync.RWMutex

	id  uint32
	ctr uint32

	numOfReplicas uint32

	appReqChs map[int]chan *shardpb.AppendRequest
	appRspCh  chan *shardpb.AppendResponse

	pbClients map[int]shardpb.ReplicaClient

	waitingAppends map[uint64]chan uint64
}

func NewClient(id uint32, replicaIPs []string) (*Client, error) {
	c := new(Client)
	c.id = id
	c.numOfReplicas = uint32(len(replicaIPs))
	c.appReqChs = make(map[int]chan *shardpb.AppendRequest)
	c.appRspCh = make(chan *shardpb.AppendResponse, 2048)
	c.pbClients = make(map[int]shardpb.ReplicaClient)
	c.waitingAppends = make(map[uint64]chan uint64)

	for i, ip := range replicaIPs {
		conn, err := grpc.Dial(ip, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		pbClient := shardpb.NewReplicaClient(conn)
		c.pbClients[i] = pbClient

		stream, err := pbClient.Append(context.Background())
		if err != nil {
			return nil, err
		}
		ch := make(chan *shardpb.AppendRequest, 2048)
		c.appReqChs[i] = ch
		go sendAppendRequests(ch, stream)
		go c.receiveAppendResponses(stream)
	}
	go c.handleAppendResponses()
	return c, nil
}

func (c *Client) Append(record string, color uint32) uint64 {
	token := c.getNewToken()
	waitingGsn := make(chan uint64, 1)
	c.mu.Lock()
	c.waitingAppends[token] = waitingGsn
	c.mu.Unlock()

	responsible := int(c.id) % len(c.appReqChs)
	for i, ch := range c.appReqChs {
		ch <- &shardpb.AppendRequest{Token: token, Color: color, Record: record, Responsible: responsible == i}
	}
	gsn := <-waitingGsn

	defer func() {
		c.mu.Lock()
		delete(c.waitingAppends, token)
		c.mu.Unlock()
	}()

	return gsn
}

func (c *Client) Read(gsn uint64, color uint32) string {
	rsp, err := c.pbClients[int(c.id)%len(c.appReqChs)].Read(context.Background(), &shardpb.ReadRequest{Gsn: gsn, Color: color})
	if err != nil {
		logrus.Fatalln(err)
	}
	return rsp.Record
}

func (c *Client) Trim(gsn uint64, color uint32) {
	for _, pbCl := range c.pbClients {
		_, err := pbCl.Trim(context.Background(), &shardpb.TrimRequest{Color: color, Gsn: gsn})
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (c *Client) handleAppendResponses() {
	numOfAppends := make(map[uint64]uint32)
	for appRsp := range c.appRspCh {
		numOfAppends[appRsp.Token]++
		if numOfAppends[appRsp.Token] == c.numOfReplicas {
			c.mu.RLock()
			c.waitingAppends[appRsp.Token] <- appRsp.Gsn
			c.mu.RUnlock()
		}
	}
}

func sendAppendRequests(ch chan *shardpb.AppendRequest, stream shardpb.Replica_AppendClient) {
	for appReq := range ch {
		err := stream.Send(appReq)
		if err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (c *Client) receiveAppendResponses(stream shardpb.Replica_AppendClient) {
	for {
		appRsp, err := stream.Recv()
		if err != nil {
			logrus.Fatalln(err)
		}
		c.appRspCh <- appRsp
	}
}

func (c *Client) getNewToken() uint64 {
	ctr := atomic.AddUint32(&c.ctr, 1)
	return (uint64(c.id) << 32) + uint64(ctr)
}
