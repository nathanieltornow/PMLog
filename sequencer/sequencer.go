package sequencer

import (
	"fmt"
	seqclient "github.com/nathanieltornow/PMLog/sequencer/seq_client"
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

const (
	heartBeatInterval = time.Second
	batchingInterval  = time.Second
)

type Sequencer struct {
	pb.UnimplementedSequencerServer

	id     uint32
	epoch  uint32
	root   bool
	leader bool
	color  uint32

	sn   uint32
	snMu sync.Mutex

	parentClient *seqclient.Client

	oReqCache   map[uint64]*pb.OrderRequest
	oReqCacheMu sync.Mutex

	oRspCs   map[uint32]chan *pb.OrderResponse
	oRspCsID uint32
	oRspCsMu sync.RWMutex

	oReqCIn chan *pb.OrderRequest
}

func NewSequencer(root bool, leader bool, color uint32) *Sequencer {
	s := new(Sequencer)
	s.leader = leader
	s.root = root
	s.color = color
	s.oReqCache = make(map[uint64]*pb.OrderRequest)
	s.oRspCs = make(map[uint32]chan *pb.OrderResponse)
	s.oReqCIn = make(chan *pb.OrderRequest, 256)
	return s
}

func (s *Sequencer) Start(IP string, parentIP string) error {
	if s.root && s.leader {
		return s.startGRPCServer(IP)
	}

	client, err := seqclient.NewClient(parentIP)
	if err != nil {
		return fmt.Errorf("failed to connect to parent: %v", err)
	}
	s.parentClient = client
	go s.handleOrderResponses()

	return s.startGRPCServer(IP)
}

func (s *Sequencer) startGRPCServer(IP string) error {
	lis, err := net.Listen("tcp", IP)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	pb.RegisterSequencerServer(server, s)

	go s.handleOrderRequests()

	go func() {
		for {
			<-time.After(10 * time.Second)
			s.snMu.Lock()
			logrus.Infoln("Sn:", s.sn)
			s.snMu.Unlock()
		}
	}()

	logrus.Infoln("starting sequencer on ", IP)
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start sequencer: %v", err)
	}
	return nil
}

func (s *Sequencer) GetOrder(stream pb.Sequencer_GetOrderServer) error {
	oRspC := make(chan *pb.OrderResponse, 256)

	s.oRspCsMu.Lock()
	id := s.oRspCsID
	s.oRspCs[s.oRspCsID] = oRspC
	s.oRspCsID++
	s.oRspCsMu.Unlock()
	go s.forwardOrderResponses(stream, oRspC)

	for {
		oReq, err := stream.Recv()
		if err != nil {
			s.oRspCsMu.Lock()
			delete(s.oRspCs, id)
			s.oRspCsMu.Unlock()
			close(oRspC)
			return err
		}
		s.oReqCIn <- oReq
	}
}

func (s *Sequencer) getAndIncSequenceNum(inc uint32) uint64 {
	s.snMu.Lock()
	defer s.snMu.Unlock()
	res := (uint64(s.epoch) << 32) + uint64(s.sn)
	s.sn += inc
	return res
}
