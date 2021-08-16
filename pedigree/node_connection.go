package pedigree

import (
	"fmt"
	"google.golang.org/grpc"
)

func newConnection(ip string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node: %v", err)
	}
	return conn, nil
}
