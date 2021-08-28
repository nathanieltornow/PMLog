package replication_client

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
)

type ReplicationClient struct {
	mu            sync.Mutex
	numOfReplicas uint32
	snToAcks      map[uint32]uint32

	comRepCs map[uint32]chan *shardpb.ReplicaMessage
	repIDCtr uint32
	repMu    sync.RWMutex

	ackC    chan *shardpb.ReplicaMessage
	comRepC chan *shardpb.ReplicaMessage
}

func NewReplicationClient() (*ReplicationClient, error) {
	repCl := new(ReplicationClient)
	repCl.ackC = make(chan *shardpb.ReplicaMessage, 1024)
	repCl.comRepC = make(chan *shardpb.ReplicaMessage, 1024)
	repCl.comRepCs = make(map[uint32]chan *shardpb.ReplicaMessage)
	go repCl.broadcastReplicationMessages()
	return nil, nil
}

func (r *ReplicationClient) AddReplica(IP string) error {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := shardpb.NewNodeClient(conn)
	stream, err := client.Replicate(context.Background())

	if err != nil {
		return fmt.Errorf("failed to start replication stream: %v", err)
	}
	r.mu.Lock()
	r.numOfReplicas++
	r.mu.Unlock()

	go r.forwardReplicationMessage(stream)
	go r.receiveReplicationMessages(stream)
	return nil
}

func (r *ReplicationClient) GetAcknowledgements() <-chan *shardpb.ReplicaMessage {
	return r.ackC
}

func (r *ReplicationClient) MakeComRepMessage() chan<- *shardpb.ReplicaMessage {
	return r.comRepC
}

func (r *ReplicationClient) receiveReplicationMessages(stream shardpb.Node_ReplicateClient) {
	for {
		in, err := stream.Recv()
		if err != nil {
			logrus.Errorf("failed to receive replica-msg: %v", err)
			return
		}

		if in.Type == shardpb.ReplicaMessage_ACK {
			r.mu.Lock()
			_, ok := r.snToAcks[in.SN]
			if !ok {
				r.snToAcks[in.SN] = 1
			}
			if r.snToAcks[in.SN] == r.numOfReplicas {
				r.ackC <- in
				delete(r.snToAcks, in.SN)
			}
			r.snToAcks[in.SN] += 1
			r.mu.Unlock()
		}
		// ignore other types of messages
	}
}

func (r *ReplicationClient) broadcastReplicationMessages() {
	for msg := range r.comRepC {
		r.repMu.RLock()
		if len(r.comRepCs) == 0 {
			if msg.Type == shardpb.ReplicaMessage_REP {
				r.ackC <- &shardpb.ReplicaMessage{Type: shardpb.ReplicaMessage_ACK, SN: msg.SN}
			}
		}
		for _, ch := range r.comRepCs {
			ch <- msg
		}
		r.repMu.RUnlock()
	}
}

func (r *ReplicationClient) forwardReplicationMessage(stream shardpb.Node_ReplicateClient) {
	comRepC := make(chan *shardpb.ReplicaMessage, 512)
	r.repMu.Lock()
	id := r.repIDCtr
	r.comRepCs[id] = comRepC
	r.repIDCtr++
	r.repMu.Unlock()

	defer func() {
		r.repMu.Lock()
		delete(r.comRepCs, id)
		r.repMu.Unlock()
		close(comRepC)
	}()

	for msg := range comRepC {
		err := stream.Send(msg)
		if err != nil {
			logrus.Errorf("failed to send rep-msg: %v", err)
			return
		}
	}
}
