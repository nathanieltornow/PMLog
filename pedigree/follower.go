package pedigree

import (
	"context"
	"fmt"
	pb "github.com/nathanieltornow/PMLog/pedigree/pedigreepb"
	"github.com/sirupsen/logrus"
	"time"
)

func (n *Node) receiveHeartbeat(stream pb.Node_HeartbeatClient) {
	defer n.selectNewLeader()
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
			n.Lock()
			n.peers = in.Peers
			n.adopters = in.Adopters
			n.Unlock()
		case <-time.After(2 * heartBeatInterval):
			return
		}

	}
}

func (n *Node) selectNewLeader() {
	if len(n.peers) == 0 {
		logrus.Fatalf("no peers")
	}
	for _, v := range n.peers {
		if v.IP == n.nodeInfo.IP {
			n.Lock()
			n.isLeader = true
			client, err := NewClient(n.adopters, n)
			if err == nil {
				n.client = client
			}
			n.peers = make([]*pb.NodeInfo, 0)
			n.adopters = make([]*pb.NodeInfo, 0)
			n.nodeIDCtr = n.nodeInfo.ID + 1
			n.Unlock()
			logrus.Infoln("Leader now")
			return
		}

		err := n.connectToLeader(v.IP)
		if err != nil {
			continue
		}
		return
	}
	logrus.Fatalln("failed to connect to any of the peers")
}

func (n *Node) connectToLeader(ip string) error {
	nodeConn, err := newConnection(ip)
	if err != nil {
		return err
	}
	client := pb.NewNodeClient(nodeConn)
	resp, err := client.Register(context.Background(), n.nodeInfo)
	if err != nil {
		return fmt.Errorf("failed to register at the leader: %v", err)
	}

	// if registered successfully, change variables
	n.Lock()
	n.nodeInfo.ID = resp.NodeID
	n.peers = resp.Structure.Peers
	n.adopters = resp.Structure.Adopters
	n.Unlock()

	stream, err := client.Heartbeat(context.Background(), &pb.Empty{})
	if err != nil {
		return err
	}
	go n.receiveHeartbeat(stream)
	return nil
}
