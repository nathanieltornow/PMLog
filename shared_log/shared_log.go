package shared_log

import (
	"context"
	pb "github.com/nathanieltornow/PMLog/shared_log/shared_logpb"
	"github.com/nathanieltornow/PMLog/shared_log/storage"
)

type SharedLog struct {
	pb.UnimplementedSharedLogServer
	log storage.Log

	newAppends chan *pb.AppendRequest

	pendingCommits map[uint64]*waitingRecord
}

func NewSharedLog(log storage.Log) (*SharedLog, error) {
	sl := new(SharedLog)
	sl.log = log
	sl.newAppends = make(chan *pb.AppendRequest, 1024)
	return sl, nil
}

func (sl *SharedLog) Append(_ context.Context, req *pb.AppendRequest) (*pb.AppendResponse, error) {
	sl.newAppends <- req

	return nil, nil
}

func (sl *SharedLog) Read(req *pb.ReadRequest, stream pb.SharedLog_ReadServer) error {
	return nil
}

func (sl *SharedLog) Trim(_ context.Context, req *pb.TrimRequest) (*pb.TrimResponse, error) {

	return nil, nil
}
