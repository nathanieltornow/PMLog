package shard

import (
	"github.com/nathanieltornow/PMLog/shard/replication_client"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/nathanieltornow/PMLog/storage"
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

	node.repClient = replication_client.NewReplicationClient()
	node.replicaOutCh = node.repClient.BroadcastReplicaMessages()
	node.ackInC = node.repClient.GetAcknowledgements()

	return node, nil
}
