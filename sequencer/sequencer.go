package sequencer

import (
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"google.golang.org/grpc"
	"sync"
	"time"
)

const (
	heartBeatInterval = time.Second
	batchingInterval  = time.Second
)

type Sequencer struct {
	pb.UnimplementedSequencerServer

	id      uint32
	root    bool
	leader  bool
	epoch   uint32
	color   uint32
	sn      uint32
	stateMu sync.Mutex

	peers    []*pb.NodeInfo
	adopters []*pb.NodeInfo
	structMu sync.Mutex

	parentConn *grpc.ClientConn
	upstream   pb.Sequencer_GetOrderClient

	oReqCache   map[uint64]*pb.OrderRequest
	oReqCacheMu sync.Mutex

	oRspCs   map[uint32]chan *pb.OrderResponse
	oRspCsMu sync.Mutex
}

func NewSequencer() (*Sequencer, error) {

	return nil, nil
}

func (s *Sequencer) GetOrder(stream pb.Sequencer_GetOrderServer) error {
	for {
		//oReq, err := stream.Recv()
		//if err != nil {
		//	return err
		//}
	}
}

func (s *Sequencer) getAndIncSequenceNum(inc uint32) uint64 {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()
	res := (uint64(s.epoch) << 32) + uint64(s.sn)
	s.sn += inc
	return res
}
