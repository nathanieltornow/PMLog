package sequencer

import (
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"time"
)

func (s *Sequencer) receiveOrderResponses() {

}

func (s *Sequencer) handleOrderRequests(oReqC chan *pb.OrderRequest) {
	s.stateMu.Lock()
	justFwd := !s.leader
	justReply := s.root && s.leader
	s.stateMu.Unlock()

	if justFwd {
		go forwardOrderRequests(s.upstream, oReqC)
		return
	}
	if justReply {
		for oReq := range oReqC {
			sn := s.getAndIncSequenceNum(oReq.NumOfRecords)
			oRsp := &pb.OrderResponse{Color: oReq.Color, NumOfRecords: oReq.NumOfRecords, Lsn: oReq.Lsn, Gsn: sn}
			s.broadcastOrderResponse(oRsp)
		}
		return
	}

	color := s.color

	batchedOReqC := make(chan *pb.OrderRequest, 256)
	go forwardOrderRequests(s.upstream, batchedOReqC)

	var newOReq *pb.OrderRequest
	for {
		select {
		case <-time.Tick(batchingInterval):
			if newOReq != nil {
				batchedOReqC <- newOReq
				newOReq = nil
			}
		case oReq := <-oReqC:
			sn := s.getAndIncSequenceNum(oReq.NumOfRecords)
			if oReq.Color == color {
				oRsp := &pb.OrderResponse{Color: oReq.Color, NumOfRecords: oReq.NumOfRecords, Lsn: oReq.Lsn, Gsn: sn}
				s.broadcastOrderResponse(oRsp)
				continue
			}

			// cache orderRequest
			s.oReqCacheMu.Lock()
			s.oReqCache[sn] = oReq
			s.oReqCacheMu.Unlock()

			if newOReq == nil {
				newOReq = &pb.OrderRequest{Lsn: sn, Color: color, NumOfRecords: oReq.NumOfRecords}
				continue
			}
			newOReq.NumOfRecords += oReq.NumOfRecords
		}
	}

}

func forwardOrderRequests(stream pb.Sequencer_GetOrderClient, oReqC chan *pb.OrderRequest) {
	for oReq := range oReqC {
		err := stream.Send(oReq)
		if err != nil {
			return
		}
	}
}

func (s *Sequencer) broadcastOrderResponse(oRsp *pb.OrderResponse) {
	s.oRspCsMu.Lock()
	for _, oRspC := range s.oRspCs {
		oRspC <- oRsp
	}
	s.oRspCsMu.Unlock()
}

func (s *Sequencer) forwardResponses(stream pb.Sequencer_GetOrderServer, oRspC chan *pb.OrderResponse) {
	for {
		s.structMu.Lock()
		resp := &pb.Response{Struct: &pb.Structure{Peers: s.peers, Adopters: s.adopters}}
		s.structMu.Unlock()
		select {
		case oReq, ok := <-oRspC:
			if !ok {
				return
			}
			resp.ORsp = oReq
			err := stream.Send(resp)
			if err != nil {
				return
			}
		case <-time.After(heartBeatInterval):
			err := stream.Send(resp)
			if err != nil {
				return
			}
		}
	}
}
