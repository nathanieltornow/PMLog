package shared_log

import (
	"context"
	pb "github.com/nathanieltornow/PMLog/shared_log/shared_logpb"
)

type SharedLog struct {
	pb.UnimplementedSharedLogServer
}

func NewSharedLog() (*SharedLog, error) {
	sl := new(SharedLog)
	return sl, nil
}

func (sl *SharedLog) Append(_ context.Context, req *pb.AppendRequest) (*pb.AppendResponse, error) {

	return nil, nil
}

func (sl *SharedLog) Read(req *pb.ReadRequest, stream pb.SharedLog_ReadServer) error {
	return nil
}

func (sl *SharedLog) Trim(_ context.Context, req *pb.TrimRequest) (*pb.TrimResponse, error) {

	return nil, nil
}
