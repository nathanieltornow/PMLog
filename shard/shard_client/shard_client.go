package shard_client

import (
	"context"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"google.golang.org/grpc"
)

type Client struct {
	pbClient shardpb.NodeClient
}

func NewClient(IP string) (*Client, error) {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	pbClient := shardpb.NewNodeClient(conn)
	client := new(Client)
	client.pbClient = pbClient
	return client, nil
}

func (c *Client) Append(record string, color uint32) (uint64, error) {
	rsp, err := c.pbClient.Append(context.Background(), &shardpb.AppendRequest{Color: color, Record: record})
	if err != nil {
		return 0, err
	}
	return rsp.GSN, nil
}
