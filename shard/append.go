package shard

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/sirupsen/logrus"
)

func (n *Node) Append(_ context.Context, request *shardpb.AppendRequest) (*shardpb.AppendResponse, error) {
	n.ctrMu.Lock()
	sn := (uint64(n.ctr) << 32) + uint64(n.ID)
	n.ctr++
	n.ctrMu.Unlock()

	inRec := &incomingRecord{record: request.Record, sn: sn, replicated: make(chan bool)}
	n.primWrC <- inRec

	// wait for the record to be replicated
	<-inRec.replicated
	err := n.primLog.Commit(sn, 0, sn)
	if err != nil {
		return nil, fmt.Errorf("failed to commit record %v", err)
	}
	return &shardpb.AppendResponse{SN: sn, GSN: sn}, nil
}

func (n *Node) storePrimaryRecords() {
	for inRec := range n.primWrC {
		err := n.primLog.Append(inRec.record, inRec.sn)
		if err != nil {
			logrus.Errorf("failed to append primary record: %v", err)
			continue
		}
		n.pendingRecordsMu.Lock()
		n.pendingRecords[inRec.sn] = inRec
		n.pendingRecordsMu.Unlock()
		n.replicaOutCh <- &shardpb.ReplicaMessage{Type: shardpb.ReplicaMessage_REP, SN: inRec.sn, Record: inRec.record}
	}
}

func (n *Node) listenForAcks() {
	for repMsg := range n.ackInC {
		n.pendingRecordsMu.RLock()
		pendRec, ok := n.pendingRecords[repMsg.SN]
		n.pendingRecordsMu.RUnlock()
		if !ok {
			logrus.Errorf("can't find pending record")
			continue
		}
		pendRec.replicated <- true
		n.replicaOutCh <- &shardpb.ReplicaMessage{Type: shardpb.ReplicaMessage_COM, SN: repMsg.SN}
	}

}
