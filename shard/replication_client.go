package shard

import (
	"context"
	"fmt"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"google.golang.org/grpc"
)

type replicationClient struct {
	stream  shardpb.Shard_ReplicateClient
	oRspC   chan *shardpb.OrderResponse
	repReqC chan *shardpb.ReplicationRequest
}

func newReplicationClient(IP string) (*replicationClient, error) {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to replica: %v", err)
	}
	client := shardpb.NewShardClient(conn)
	stream, err := client.Replicate(context.Background())
	if err != nil {
		return nil, err
	}
	repCl := new(replicationClient)
	repCl.stream = stream
	repCl.oRspC = make(chan *shardpb.OrderResponse, 1024)
	repCl.repReqC = make(chan *shardpb.ReplicationRequest, 1024)
	return repCl, nil
}

func (r *replicationClient) MakeReplicationRequest() chan<- *shardpb.ReplicationRequest {
	return r.repReqC
}

func (r *replicationClient) GetOrderResponses() <-chan *shardpb.OrderResponse {
	return r.oRspC
}

func (r *replicationClient) Stop() error {
	err := r.stream.CloseSend()
	if err != nil {
		return err
	}
	close(r.oRspC)
	close(r.repReqC)
	return nil
}

func (r *replicationClient) receiveORsps() {
	for {
		rsp, err := r.stream.Recv()
		if err != nil {
			return
		}
		r.oRspC <- rsp
	}
}

func (r *replicationClient) sendReplicationRequests() {
	for rep := range r.repReqC {
		err := r.stream.Send(rep)
		if err != nil {
			return
		}
	}
}
