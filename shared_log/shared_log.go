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

func (sl *SharedLog) Start(ipAddr, nodeIP, orderIP string, peerIPs []string, interval time.Duration) error {
	lis, err := net.Listen("tcp", ipAddr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterSharedLogServer(server, sl)

	node, err := app_node.NewNode(sl.id, sl.color, app_node.WithBatchingInterval(interval))
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
	start := time.Now()
	localToken := sl.node.MakeCommitRequest(&frame.CommitRequest{Color: req.Color, Content: req.Record})
	fmt.Println("one", time.Since(start))
	sl.pendingAppendsMu.Lock()
	sl.pendingAppends[localToken] = waitingGsn
	sl.pendingAppendsMu.Unlock()
	// wait for the global-sequence number
	s2 := time.Now()
	gsn := <-waitingGsn
	fmt.Println("two", time.Since(s2))
	return &pb.AppendResponse{Gsn: gsn}, nil
}

func (sl *SharedLog) Read(_ context.Context, req *pb.ReadRequest) (*pb.ReadResponse, error) {
	record, err := sl.log.Read(req.Gsn)
	if err != nil {
		return nil, err
	}
	if record == "" {
		return nil, fmt.Errorf("failed to find record with gsn %v", req.Gsn)
	}
	return &pb.ReadResponse{Gsn: req.Gsn, Record: record}, nil
}

func (sl *SharedLog) Trim(_ context.Context, req *pb.TrimRequest) (*pb.TrimResponse, error) {
	if err := sl.log.Trim(req.Gsn); err != nil {
		return nil, err
	}
	return &pb.TrimResponse{}, nil
}
