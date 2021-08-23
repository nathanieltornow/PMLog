package sequencer

import (
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
)

func (s *Sequencer) handleOrderResponses() {
	justFwd := !s.leader

	oRspC := s.parentClient.GetOrderResponses()

	if justFwd {
		for oRsp := range oRspC {
			go s.broadcastOrderResponse(oRsp)
		}
		return
	}

	color := s.color
	for oRsp := range oRspC {
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
				go s.broadcastOrderResponse(newORsp)

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

	oReqC := s.parentClient.MakeOrderRequests()

	if justFwd {
		for oReq := range s.oReqCIn {
			oReqC <- oReq
		}
		return
	}

	for oReq := range s.oReqCIn {
		sn := s.getAndIncSequenceNum(oReq.NumOfRecords)

		if color == oReq.Color {
			oRsp := &pb.OrderResponse{
				Lsn:          oReq.Lsn,
				Gsn:          sn,
				NumOfRecords: oReq.NumOfRecords,
				Color:        color,
				OriginColor:  oReq.OriginColor,
			}
			go s.broadcastOrderResponse(oRsp)
			continue
		}

		// cache orderRequest
		s.oReqCacheMu.Lock()
		s.oReqCache[sn] = oReq
		s.oReqCacheMu.Unlock()

		oReqC <- oReq
	}
}

func (s *Sequencer) broadcastOrderResponse(oRsp *pb.OrderResponse) {
	s.oRspCsMu.RLock()
	for _, oRspC := range s.oRspCs {
		oRspC <- oRsp
	}
	s.oRspCsMu.RUnlock()
}

func (s *Sequencer) forwardOrderResponses(stream pb.Sequencer_GetOrderServer, oRspC chan *pb.OrderResponse) {
	for oRsp := range oRspC {
		err := stream.Send(oRsp)
		if err != nil {
			return
		}
	}
}
