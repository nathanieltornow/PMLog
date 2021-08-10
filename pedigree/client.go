package pedigree

import (
	"context"
	"fmt"
	pb "github.com/nathanieltornow/PMLog/pedigree/pedigreepb"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Client struct {
	sync.Mutex

	parentInfo *pb.NodeInfo
	parentConn *nodeConnection

	adopters []*pb.NodeInfo
}

func NewClient(adopters []*pb.NodeInfo) (*Client, error) {
	if len(adopters) == 0 {
		return nil, fmt.Errorf("no adopters")
	}
	client := new(Client)
	client.adopters = adopters
	client.adopt()
	return client, nil
}

func (c *Client) GetParentAppIP() string {
	c.Lock()
	defer c.Unlock()
	return c.parentInfo.IP + ":" + c.parentInfo.AppPort
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
			c.Lock()
			c.adopters = in.Peers
			c.Unlock()
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
		nodeConn, err := newNodeConnection(v)
		if err != nil {
			continue
		}
		c.Lock()
		c.parentInfo = v
		c.parentConn = nodeConn
		c.Unlock()
		stream, err := c.parentConn.client.Heartbeat(context.Background(), &pb.Empty{})
		if err != nil {
			logrus.Errorf("failed to start heartbeat: %v", err)
			continue
		}
		logrus.Infoln("Adopted by:", c.parentInfo)
		go c.receiveHeartbeat(stream)
		return
	}
	logrus.Fatalln("failed to connect to any of the adopters")
}
