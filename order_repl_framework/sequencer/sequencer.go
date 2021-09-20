package sequencer

import (
	"fmt"
	seq_client "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/client"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"time"
)

type Sequencer struct {
	sequencerpb.UnimplementedSequencerServer

	id     uint32
	epoch  uint32
	root   bool
	leader bool
	color  uint32

	sn   uint32
	snMu sync.Mutex

	parentClient *seq_client.Client

	oRspCs   map[uint64]chan *sequencerpb.OrderResponse
	oRspCsID uint32
	oRspCsMu sync.RWMutex

	broadcastC chan *sequencerpb.OrderResponse

	oReqCIn chan *sequencerpb.OrderRequest
}

func NewSequencer(root bool, leader bool, color uint32) *Sequencer {
	s := new(Sequencer)
	s.leader = leader
	s.root = root
	s.color = color
	s.oRspCs = make(map[uint64]chan *sequencerpb.OrderResponse)
	s.oReqCIn = make(chan *sequencerpb.OrderRequest, 1024)
	s.broadcastC = make(chan *sequencerpb.OrderResponse, 1024)
	return s
}

func (s *Sequencer) Start(IP string, parentIP string) error {
	if s.root && s.leader {
		return s.startGRPCServer(IP)
	}

	client, err := seq_client.NewClient(parentIP, s.color)
	if err != nil {
		return fmt.Errorf("failed to connect to parent: %v", err)
	}
	s.parentClient = client

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
	if !s.root {
		go s.handleOrderResponses()
	}
	go s.broadcastOrderResponses()

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

func (s *Sequencer) GetOrder(stream sequencerpb.Sequencer_GetOrderServer) error {
	oRspC := make(chan *sequencerpb.OrderResponse, 256)
	var id uint64
	first := true

	go s.forwardOrderResponses(stream, oRspC)

	for {
		batchedOReq, err := stream.Recv()
		if err != nil {
			s.oRspCsMu.Lock()
			delete(s.oRspCs, id)
			s.oRspCsMu.Unlock()
			close(oRspC)
			return err
		}

		if first {
			s.oRspCsMu.Lock()
			s.oRspCsID++
			id = (uint64(s.oRspCsID) << 32) + uint64(batchedOReq.OReqs[0].OriginColor)
			s.oRspCs[id] = oRspC
			s.oRspCsMu.Unlock()
			first = false
		}

		for _, oReq := range batchedOReq.OReqs {
			s.oReqCIn <- oReq
		}
	}
}

func (s *Sequencer) getAndIncSequenceNum(inc uint32) uint64 {
	s.snMu.Lock()
	defer s.snMu.Unlock()
	res := (uint64(s.epoch) << 32) + uint64(s.sn)
	s.sn += inc
	return res
}
