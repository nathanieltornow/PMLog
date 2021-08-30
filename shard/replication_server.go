package shard

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/sirupsen/logrus"
)

func (n *Node) Register(_ context.Context, request *shardpb.RegisterRequest) (*shardpb.OK, error) {
	err := n.repClient.AddReplica(request.IP)
	if err != nil {
		return nil, fmt.Errorf("failed to register replica: %v", err)
	}
	return &shardpb.OK{}, nil
}

func (n *Node) Replicate(stream shardpb.Node_ReplicateServer) error {
	go n.replyAcks(stream)
	for {
		msg, err := stream.Recv()
		if err != nil {
			logrus.Errorf("failed to receive rep-request: %v", err)
			return err
		}
		if msg.Type == shardpb.ReplicaMessage_REP {
			n.repInC <- msg
		}
		if msg.Type == shardpb.ReplicaMessage_COM {
			n.comInC <- msg
		}
	}
}

func (n *Node) handleCommitMsgs() {
	for repMsg := range n.comInC {
		err := n.secLog.Commit(repMsg.SN, 0, repMsg.SN)
		if err != nil {
			logrus.Errorf("failed to commit in secLog: %v", err)
		}
	}
}

func (n *Node) handleRepMsgs() {
	for repMsg := range n.repInC {
		coorCtr := uint32(repMsg.SN >> 32)
		n.ctrMu.Lock()
		if coorCtr > n.ctr {
			n.ctr = coorCtr + 1
		}
		n.ctrMu.Unlock()
		err := n.secLog.Append(repMsg.Record, repMsg.SN)
		if err != nil {
			logrus.Errorf("failed to append secLog: %v", err)
		}
		n.ackOutC <- repMsg
	}
}

func (n *Node) replyAcks(stream shardpb.Node_ReplicateServer) {
	for repMsg := range n.ackOutC {
		err := stream.Send(repMsg)
		if err != nil {
			logrus.Errorf("failed to send ackMsg")
		}
	}
}
