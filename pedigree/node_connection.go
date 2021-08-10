package pedigree

import (
	"fmt"
	pb "github.com/nathanieltornow/PMLog/pedigree/pedigreepb"
	"google.golang.org/grpc"
)

type nodeConnection struct {
	nodeInfo *pb.NodeInfo
	conn     *grpc.ClientConn
	client   pb.NodeClient
}

func newNodeConnection(info *pb.NodeInfo) (*nodeConnection, error) {
	conn, err := grpc.Dial(info.IP+":"+info.Port, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node: %v", err)
	}
	client := pb.NewNodeClient(conn)
	return &nodeConnection{nodeInfo: info, conn: conn, client: client}, nil
}
