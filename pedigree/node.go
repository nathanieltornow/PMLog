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

	isLeader bool
	epoch    uint32
	nodeInfo *pb.NodeInfo

	app App

	peers    []*pb.NodeInfo
	adopters []*pb.NodeInfo

	// client to a overlying group, only useful for leader nodes
	client *Client
	// counter for giving follower nodes a distinct and increasing ID
	nodeIDCtr uint32
}

type NodeState struct {
	IsLeader   bool
	Epoch      uint32
	ParentConn *grpc.ClientConn
}

type App interface {
	Start(listener net.Listener, parentConn *grpc.ClientConn) error
	UpdateState(state NodeState) error
}

// NewNode creates a new node, considering if it is a leader or a follower
func NewNode(isLeader bool, app App) (*Node, error) {
	n := new(Node)
	n.isLeader = isLeader
	n.peers = make([]*pb.NodeInfo, 0)
	n.adopters = make([]*pb.NodeInfo, 0)
	n.app = app
	if isLeader {
		n.nodeIDCtr = 1
	}
	return n, nil
}

// Start starts a new node, considering if it is a leader or a follower
func (n *Node) Start(nodeInfo, parentNodeInfo *pb.NodeInfo) error {
	n.nodeInfo = nodeInfo
	if parentNodeInfo == nil {
		if !n.isLeader {
			return fmt.Errorf("failed to start: has to be leader if there is no parent")
		}
	} else {
		if n.isLeader {
			client, err := NewClient([]*pb.NodeInfo{parentNodeInfo}, n)
			if err != nil {
				return err
			}
			n.client = client
		} else {
			err := n.connectToLeader()
			if err != nil {
				return err
			}
		}
	}
	lis, err := net.Listen("tcp", nodeInfo.IP)
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

// Register is the RPC method that registers a follower at the leader.
// If the node is the leader, it returns successfully with the ID of the requesting node and the structure of the tree
func (n *Node) Register(_ context.Context, req *pb.NodeInfo) (*pb.RegisterResponse, error) {
	if !n.isLeader {
		return nil, fmt.Errorf("failed to register: node is non-leader")
	}
	var adopters []*pb.NodeInfo
	if n.client != nil {
		adopters = n.client.adopters
	}
	n.Lock()
	n.peers = append(n.peers, &pb.NodeInfo{ID: n.nodeIDCtr, IP: req.IP})
	structure := pb.StructureUpdate{Epoch: n.nodeInfo.ID, Peers: n.peers, Adopters: adopters}
	resp := pb.RegisterResponse{NodeID: n.nodeIDCtr, Structure: &structure}
	n.nodeIDCtr++
	n.Unlock()
	return &resp, nil
}

// Heartbeat is the RPC method that followers and clients subscribe to in order to see if the parent is running healthy
// Heartbeat also gives the followers and clients updates of the structure
func (n *Node) Heartbeat(_ *pb.Empty, stream pb.Node_HeartbeatServer) error {
	ticker := time.Tick(heartBeatInterval)
	for range ticker {
		var adopters []*pb.NodeInfo
		if n.client != nil {
			adopters = n.client.adopters
		}
		n.Lock()
		update := pb.StructureUpdate{Epoch: n.nodeInfo.ID, Peers: n.peers, Adopters: adopters}
		n.Unlock()
		err := stream.Send(&update)
		if err != nil {
			return fmt.Errorf("failed to send heartbeat: %v", err)
		}
	}
	return nil
}

func (n *Node) UpdateParentConnection(parentConn *grpc.ClientConn) {

}
