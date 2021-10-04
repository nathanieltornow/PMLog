package shared_log

import (
	"context"
	"fmt"
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node"
	pb "github.com/nathanieltornow/PMLog/shared_log/shared_logpb"
	"github.com/nathanieltornow/PMLog/shared_log/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type SharedLog struct {
	pb.UnimplementedSharedLogServer
	log storage.Log

	id    uint32
	color uint32

	node *app_node.Node

	pendingAppends   map[uint64]chan uint64
	pendingAppendsMu sync.Mutex
}

func NewSharedLog(log storage.Log, id, color uint32) (*SharedLog, error) {
	sl := new(SharedLog)
	sl.id = id
	sl.color = color
	sl.log = log
	sl.pendingAppends = make(map[uint64]chan uint64)
	return sl, nil
}

func (sl *SharedLog) Start(ipAddr, nodeIP, orderIP string, peerIPs []string) error {
	lis, err := net.Listen("tcp", ipAddr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterSharedLogServer(server, sl)

	node, err := app_node.NewNode(sl.id, sl.color)
	if err != nil {
		return err
	}
	sl.node = node
	node.RegisterApp(sl)
	go node.Start(nodeIP, peerIPs, orderIP)
	time.Sleep(time.Second)
	logrus.Infoln("Starting Shared log")
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start shared log: %v", err)
	}
	return nil
}

func (sl *SharedLog) Append(_ context.Context, req *pb.AppendRequest) (*pb.AppendResponse, error) {
	waitingGsn := make(chan uint64, 1)
	localToken := sl.node.MakeCommitRequest(&frame.CommitRequest{Color: req.Color, Content: req.Record})
	sl.pendingAppendsMu.Lock()
	sl.pendingAppends[localToken] = waitingGsn
	sl.pendingAppendsMu.Unlock()
	// wait for the global-sequence number
	gsn := <-waitingGsn
	return &pb.AppendResponse{Gsn: gsn}, nil
}

func (sl *SharedLog) Read(req *pb.ReadRequest, stream pb.SharedLog_ReadServer) error {
	return nil
}

func (sl *SharedLog) Trim(_ context.Context, req *pb.TrimRequest) (*pb.TrimResponse, error) {

	return nil, nil
}
