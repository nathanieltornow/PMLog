package shard

import (
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/replication_client"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/nathanieltornow/PMLog/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type incomingRecord struct {
	sn         uint64
	record     string
	replicated chan bool
}

type Node struct {
	shardpb.UnimplementedNodeServer

	ID    uint32
	color uint32

	ctr   uint32
	ctrMu sync.Mutex

	repClient *replication_client.ReplicationClient

	primLog storage.Log
	secLog  storage.Log

	primWrC chan *incomingRecord
	secWrC  chan *incomingRecord

	pendingRecords   map[uint64]*incomingRecord
	pendingRecordsMu sync.RWMutex

	comInC  chan *shardpb.ReplicaMessage
	repInC  chan *shardpb.ReplicaMessage
	ackOutC chan *shardpb.ReplicaMessage

	ackInC       <-chan *shardpb.ReplicaMessage
	replicaOutCh chan<- *shardpb.ReplicaMessage
}

func NewNode(primLog, secLog storage.Log) (*Node, error) {
	node := new(Node)
	node.primLog = primLog
	node.secLog = secLog
	node.comInC = make(chan *shardpb.ReplicaMessage, 1024)
	node.repInC = make(chan *shardpb.ReplicaMessage, 1024)
	node.ackOutC = make(chan *shardpb.ReplicaMessage, 1024)
	node.pendingRecords = make(map[uint64]*incomingRecord)

	return node, nil
}

func (n *Node) Start(IP string, replicaIPs []string) error {
	n.repClient = replication_client.NewReplicationClient()
	n.replicaOutCh = n.repClient.BroadcastReplicaMessages()
	n.ackInC = n.repClient.GetAcknowledgements()

	for _, rep := range replicaIPs {
		err := n.repClient.AddReplica(IP, rep)
		if err != nil {
			return fmt.Errorf("failed to add replica: %v", err)
		}
	}

	go n.storePrimaryRecords()
	go n.listenForAcks()

	go n.handleCommitMsgs()
	go n.handleRepMsgs()

	lis, err := net.Listen("tcp", IP)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	shardpb.RegisterNodeServer(server, n)
	logrus.Infoln("starting node on ", IP)
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start sequencer: %v", err)
	}
	return nil
}
