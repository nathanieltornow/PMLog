package shard

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/replication_client"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/sirupsen/logrus"
)

type Node struct {
	shardpb.UnimplementedNodeServer

	repClient replication_client.ReplicationClient

	comC chan *shardpb.ReplicaMessage
	repC chan *shardpb.ReplicaMessage
}

func NewNode() (*Node, error) {

	return nil, nil
}

func (n *Node) Append(_ context.Context, request *shardpb.AppendRequest) (*shardpb.AppendResponse, error) {
	return nil, nil
}

func (n *Node) Replicate(stream shardpb.Node_ReplicateServer) error {
	for {
		msg, err := stream.Recv()
		if err != nil {
			logrus.Errorf("failed to receive rep-request: %v", err)
			return err
		}
		if msg.Type == shardpb.ReplicaMessage_REP {
			n.repC <- msg
		}
		if msg.Type == shardpb.ReplicaMessage_COM {
			n.comC <- msg
		}
	}
}

func (n *Node) Register(_ context.Context, request *shardpb.RegisterRequest) (*shardpb.OK, error) {
	err := n.repClient.AddReplica(request.IP)
	if err != nil {
		return nil, fmt.Errorf("failed to register replica: %v", err)
	}
	return &shardpb.OK{}, nil
}

func (n *Node) handleCommitMsgs() {

}

func (n *Node) handleRepMsgs() {

}
