package pedigree

import (
	"context"
	"fmt"
	pb "github.com/nathanieltornow/PMLog/pedigree/pedigreepb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"time"
)

type Client struct {
	mu        sync.Mutex
	adopters  []*pb.NodeInfo
	clientApp ClientApp
}

type ClientApp interface {
	UpdateParentConnection(conn *grpc.ClientConn)
}

func NewClient(adopters []*pb.NodeInfo, clientApp ClientApp) (*Client, error) {
	if len(adopters) == 0 {
		return nil, fmt.Errorf("failed to create new client: no adopters")
	}
	client := new(Client)
	client.adopters = adopters
	client.clientApp = clientApp
	client.adopt()
	return client, nil
}

func (c *Client) receiveHeartbeat(stream pb.Node_HeartbeatClient) {
	defer c.adopt()
	for {
		inC := make(chan *pb.StructureUpdate)
		errC := make(chan error)

		go func() {
			in, err := stream.Recv()
			inC <- in
			errC <- err
		}()

		select {
		case in := <-inC:
			err := <-errC
			if err != nil {
				return
			}
			c.mu.Lock()
			c.adopters = in.Peers
			c.mu.Unlock()
		case <-time.After(2 * heartBeatInterval):
			return
		}

	}
}

func (c *Client) adopt() {
	if len(c.adopters) == 0 {
		logrus.Fatalf("no adopters")
	}
	for _, v := range c.adopters {
		conn, err := newConnection(v.IP)
		if err != nil {
			continue
		}

		client := pb.NewNodeClient(conn)
		stream, err := client.Heartbeat(context.Background(), &pb.Empty{})
		if err != nil {
			logrus.Errorf("failed to start heartbeat: %v", err)
			continue
		}

		logrus.Infoln("Adopted by:", v)
		c.clientApp.UpdateParentConnection(conn)
		go c.receiveHeartbeat(stream)
		return
	}
	logrus.Fatalln("failed to connect to any of the adopters")
}
