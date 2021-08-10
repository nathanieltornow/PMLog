package pedigree

import (
	"context"
	"fmt"
	pb "github.com/nathanieltornow/PMLog/pedigree/pedigreepb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const (
	heartBeatInterval = time.Second
)

type Node struct {
	sync.Mutex
	pb.UnimplementedNodeServer

	nodeInfo *pb.NodeInfo

	isLeader bool
	epoch    uint32

	peers    []*pb.NodeInfo
	adopters []*pb.NodeInfo

	// as follower
	leaderConn *nodeConnection
	leaderInfo *pb.NodeInfo

	// as leader
	client *Client

	nodeIDCtr uint32
}

func NewNode(isLeader bool) (*Node, error) {
	n := new(Node)
	n.isLeader = isLeader
	n.peers = make([]*pb.NodeInfo, 0)
	n.adopters = make([]*pb.NodeInfo, 0)
	if isLeader {
		n.nodeIDCtr = 1
	}
	return n, nil
}

func (n *Node) Start(nodeInfo, parentNodeInfo *pb.NodeInfo) error {
	n.nodeInfo = nodeInfo
	if parentNodeInfo == nil {
		if !n.isLeader {
			return fmt.Errorf("failed to start: has to be leader if there is no parent")
		}
	} else {
		if n.isLeader {
			client, err := NewClient([]*pb.NodeInfo{parentNodeInfo})
			if err != nil {
				return err
			}
			n.client = client
		} else {
			n.leaderInfo = parentNodeInfo
			err := n.connectToLeader()
			if err != nil {
				return err
			}
		}

	}
	lis, err := net.Listen("tcp", nodeInfo.IP+":"+nodeInfo.Port)
	if err != nil {
		return fmt.Errorf("failed to start: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNodeServer(s, n)
	logrus.Infoln("Starting node")
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("failed to start nodeInfo: %v", err)
	}
	return nil
}

func (n *Node) Register(_ context.Context, req *pb.NodeInfo) (*pb.RegisterResponse, error) {
	if !n.isLeader {
		return nil, fmt.Errorf("failed to register: node is non-leader")
	}
	var adopters []*pb.NodeInfo
	if n.client != nil {
		adopters = n.client.adopters
	}
	n.Lock()
	n.peers = append(n.peers, &pb.NodeInfo{ID: n.nodeIDCtr, IP: req.IP, AppPort: req.AppPort, Port: req.Port})
	structure := pb.StructureUpdate{Epoch: n.epoch, Peers: n.peers, Adopters: adopters}
	resp := pb.RegisterResponse{NodeID: n.nodeIDCtr, Structure: &structure}
	n.nodeIDCtr++
	n.Unlock()
	return &resp, nil
}

func (n *Node) Heartbeat(_ *pb.Empty, stream pb.Node_HeartbeatServer) error {
	ticker := time.Tick(heartBeatInterval)
	for range ticker {
		var adopters []*pb.NodeInfo
		if n.client != nil {
			adopters = n.client.adopters
		}
		n.Lock()
		update := pb.StructureUpdate{Epoch: n.epoch, Peers: n.peers, Adopters: adopters}
		n.Unlock()
		err := stream.Send(&update)
		if err != nil {
			return fmt.Errorf("failed to send heartbeat: %v", err)
		}
	}
	return nil
}
