package shared_log

import (
	"context"
	"fmt"
	pb "github.com/nathanieltornow/PMLog/shared_log/shared_logpb"
	"github.com/nathanieltornow/PMLog/shared_log/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"sync/atomic"
)

type SharedLog struct {
	pb.UnimplementedSharedLogServer
	log storage.Log

	tokenCtr         uint64
	pendingAppends   map[uint64]chan uint64
	pendingAppendsMu sync.Mutex

	localToFindToken   map[uint64]uint64
	localToFindTokenMu sync.Mutex

	newAppends chan *newRecord
}

func NewSharedLog(log storage.Log) (*SharedLog, error) {
	sl := new(SharedLog)
	sl.log = log
	sl.pendingAppends = make(map[uint64]chan uint64)
	sl.newAppends = make(chan *newRecord, 1024)
	sl.localToFindToken = make(map[uint64]uint64)
	return sl, nil
}

func (sl *SharedLog) Start(ipAddr string) error {
	lis, err := net.Listen("tcp", ipAddr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	pb.RegisterSharedLogServer(server, sl)

	logrus.Infoln("Starting Shared log")
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start shared log: %v", err)
	}
	return nil
}

func (sl *SharedLog) Append(_ context.Context, req *pb.AppendRequest) (*pb.AppendResponse, error) {
	newToken := atomic.AddUint64(&sl.tokenCtr, 1)
	waitingGsn := make(chan uint64)
	sl.pendingAppendsMu.Lock()
	sl.pendingAppends[newToken] = waitingGsn
	sl.pendingAppendsMu.Unlock()

	sl.newAppends <- &newRecord{findToken: newToken, color: req.Color, gsn: waitingGsn, record: req.Record}

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
