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
		if v.IP == n.nodeInfo.IP && v.Port == n.nodeInfo.Port {
			n.Lock()
			n.isLeader = true
			client, err := NewClient(n.adopters)
			if err == nil {
				n.client = client
			}
			n.peers = make([]*pb.NodeInfo, 0)
			n.adopters = make([]*pb.NodeInfo, 0)
			n.leaderConn = nil
			n.leaderInfo = nil
			n.nodeIDCtr = n.epoch + 1
			n.Unlock()
			logrus.Infoln("Leader now")
			return
		}

		n.Lock()
		n.leaderInfo = v
		n.Unlock()
		err := n.connectToLeader()
		if err != nil {
			continue
		}
		logrus.Infoln("Connected to leader:", n.leaderConn.nodeInfo)
		return
	}
	logrus.Fatalln("failed to connect to any of the peers")
}

func (n *Node) connectToLeader() error {
	nodeConn, err := newNodeConnection(n.leaderInfo)
	if err != nil {
		return err
	}
	client := pb.NewNodeClient(nodeConn.conn)
	resp, err := client.Register(context.Background(), n.nodeInfo)
	if err != nil {
		return fmt.Errorf("failed to register at the leader: %v", err)
	}

	// if registered successfully, change variables
	n.Lock()
	n.epoch = resp.NodeID
	n.peers = resp.Structure.Peers
	n.adopters = resp.Structure.Adopters
	n.leaderConn = nodeConn
	n.Unlock()

	stream, err := nodeConn.client.Heartbeat(context.Background(), &pb.Empty{})
	if err != nil {
		return err
	}
	go n.receiveHeartbeat(stream)
	return nil
}
