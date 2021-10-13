package sequencer

import (
	"fmt"
	"github.com/nathanieltornow/PMLog/sequencer/client"
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
)

type Sequencer struct {
	sequencerpb.UnimplementedSequencerServer

	id    uint32
	epoch uint32
	root  bool
	color uint32

	sn   uint32
	snMu sync.Mutex

	parentClient *client.Client

	//localSeqs   map[uint32]*local_sequencer.LocalSequencer
	//localSeqsMu sync.RWMutex

	broadcastCh chan *sequencerpb.OrderResponse

	oRspCs   map[idColorTuple]chan *sequencerpb.OrderResponse
	oRspCsID uint32
	oRspCsMu sync.RWMutex

	oReqCIn chan *sequencerpb.OrderRequest
}

func NewSequencer(root bool, color uint32) *Sequencer {
	s := new(Sequencer)
	s.root = root
	s.color = color
	s.oRspCs = make(map[idColorTuple]chan *sequencerpb.OrderResponse)
	s.oReqCIn = make(chan *sequencerpb.OrderRequest, 1024)
	//s.localSeqs = make(map[uint32]*local_sequencer.LocalSequencer)
	s.broadcastCh = make(chan *sequencerpb.OrderResponse, 1024)
	return s
}

func (s *Sequencer) Start(IP string, parentIP string) error {
	if s.root {
		return s.startGRPCServer(IP)
	}

	cl, err := client.NewClient(parentIP, s.color)
	if err != nil {
		return fmt.Errorf("failed to connect to parent: %v", err)
	}
	s.parentClient = cl
	go s.handleOrderResponses()

	return s.startGRPCServer(IP)
}

func (s *Sequencer) startGRPCServer(IP string) error {
	lis, err := net.Listen("tcp", IP)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	sequencerpb.RegisterSequencerServer(server, s)

	go s.handleOrderRequests()
	go s.broadcastOrderResponses()

	logrus.Infoln("starting sequencer on ", IP)
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start sequencer: %v", err)
	}
	return nil
}

func (s *Sequencer) GetOrder(stream sequencerpb.Sequencer_GetOrderServer) error {
	oRspC := make(chan *sequencerpb.OrderResponse, 256)

	first := true
	var ict idColorTuple

	go s.forwardOrderResponses(stream, oRspC)
	for {
		oReq, err := stream.Recv()
		if first {
			s.oRspCsMu.Lock()
			ict = idColorTuple{id: s.oRspCsID, color: oReq.OriginColor}
			s.oRspCs[ict] = oRspC
			s.oRspCsID++
			s.oRspCsMu.Unlock()
			first = false
		}
		if err != nil {
			s.oRspCsMu.Lock()
			delete(s.oRspCs, ict)
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

type idColorTuple struct {
	id    uint32
	color uint32
}
