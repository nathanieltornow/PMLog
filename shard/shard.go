package shard

import (
	"context"
	"fmt"
	shardpb "github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/nathanieltornow/PMLog/storage"
	"sync"
)

type Shard struct {
	shardpb.UnimplementedShardServer

	head    bool
	tail    bool
	stateMu sync.Mutex

	log storage.Log
}

func NewShard(log storage.Log, head, tail bool) (*Shard, error) {
	shard := new(Shard)
	shard.log = log
	return shard, nil
}

func (s *Shard) Append(_ context.Context, req *shardpb.AppendRequest) (*shardpb.OrderResponse, error) {
	s.stateMu.Lock()
	head := s.head
	s.stateMu.Unlock()
	if !head {
		return nil, fmt.Errorf("shard isn't head, no append possible")
	}

	return nil, nil
}

func (s *Shard) Replicate(stream shardpb.Shard_ReplicateServer) error {

	return nil
}
