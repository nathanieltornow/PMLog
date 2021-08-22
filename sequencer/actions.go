package sequencer

import (
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
)

func (s *Sequencer) handleOrderResponses() {
	justFwd := !s.leader

	defer logrus.Infoln("Returning from handleOrderResponses")

	if justFwd {
		for oRsp := range s.oRspCIn {
			s.broadcastOrderResponse(oRsp)
		}
		return
	}

	color := s.color
	for oRsp := range s.oRspCIn {
		if oRsp.OriginColor == color {
			i := uint32(0)
			for i < oRsp.NumOfRecords {
				// get orderRequests from cache
				s.oReqCacheMu.Lock()
				oReq, ok := s.oReqCache[oRsp.Lsn]
				s.oReqCacheMu.Unlock()
				if !ok {
					logrus.Fatalln("failed to get cached OrderRequest")
				}

				newORsp := &pb.OrderResponse{
					Lsn:          oReq.Lsn,
					NumOfRecords: oReq.NumOfRecords,
					Gsn:          oRsp.Gsn,
					OriginColor:  oReq.OriginColor,
					Color:        oRsp.Color,
				}
				s.broadcastOrderResponse(newORsp)

				i += oReq.NumOfRecords
			}

			continue
		}
		s.broadcastOrderResponse(oRsp)
	}
}

func (s *Sequencer) handleOrderRequests() {
	justFwd := !s.leader
	justReply := s.root && s.leader

	defer logrus.Infoln("Returning from handleOrderRequests")

	if justReply {
		// in case the sequencer is the root, it will just immediately return with an OrderResponse
		for oReq := range s.oReqCIn {
			sn := s.getAndIncSequenceNum(oReq.NumOfRecords)
			oRsp := &pb.OrderResponse{
				Lsn:          oReq.Lsn,
				Gsn:          sn,
				NumOfRecords: oReq.NumOfRecords,
				Color:        oReq.Color,
				OriginColor:  oReq.OriginColor,
			}
			s.broadcastOrderResponse(oRsp)
		}
		return
	}

	color := s.color

	batchedOReqC := make(chan *pb.OrderRequest, 256)
	go s.forwardOrderRequests(batchedOReqC)
	defer close(batchedOReqC)

	if justFwd {
		s.forwardOrderRequests(s.oReqCIn)
		return
	}

	for oReq := range s.oReqCIn {
		sn := s.getAndIncSequenceNum(oReq.NumOfRecords)
		oRsp := &pb.OrderResponse{
			Lsn:          oReq.Lsn,
			Gsn:          sn,
			NumOfRecords: oReq.NumOfRecords,
			Color:        color,
			OriginColor:  oReq.OriginColor,
		}
		s.broadcastOrderResponse(oRsp)

		if color == oReq.Color {
			continue
		}

		// cache orderRequest
		s.oReqCacheMu.Lock()
		s.oReqCache[sn] = oReq
		s.oReqCacheMu.Unlock()

		batchedOReqC <- oReq
	}
}

func (s *Sequencer) broadcastOrderResponse(oRsp *pb.OrderResponse) {
	s.oRspCsMu.RLock()
	for _, oRspC := range s.oRspCs {
		oRspC <- oRsp
	}
	s.oRspCsMu.RUnlock()
}

func (s *Sequencer) forwardOrderRequests(oReqC chan *pb.OrderRequest) {
	defer logrus.Infoln("Returning from forwardOrderRequests")
	for oReq := range oReqC {
		err := s.upstream.Send(oReq)
		if err != nil {
			return
		}
	}
}

func (s *Sequencer) forwardOrderResponses(stream pb.Sequencer_GetOrderServer, oRspC chan *pb.OrderResponse) {
	defer logrus.Infoln("Returning from OrderResponses")

	for oRsp := range oRspC {
		err := stream.Send(oRsp)
		if err != nil {
			return
		}
	}
}
