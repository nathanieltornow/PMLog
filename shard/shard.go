package shard

import "github.com/nathanieltornow/PMLog/storage"

type Shard struct {
	log storage.Log
}

func NewShard(log storage.Log) (*Shard, error) {
	shard := new(Shard)
	shard.log = log
	return shard, nil
}
