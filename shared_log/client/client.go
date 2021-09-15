package client

import (
	"context"
	pb "github.com/nathanieltornow/PMLog/shared_log/shared_logpb"
	"google.golang.org/grpc"
)

type Client struct {
	pbClient pb.SharedLogClient
}

func NewClient(IP string) (*Client, error) {
	cl := new(Client)
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	cl.pbClient = pb.NewSharedLogClient(conn)
	return cl, nil
}

func (cl *Client) Append(color uint32, record string) (uint64, error) {
	req := &pb.AppendRequest{Record: record, Color: color}
	resp, err := cl.pbClient.Append(context.Background(), req)
	if err != nil {
		return 0, err
	}
	return resp.Gsn, nil
}
