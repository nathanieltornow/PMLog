package sequencer

import (
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
)

func (s *Sequencer) handleOrderResponses() {
	for {
		oRsp := s.parentClient.GetNextOrderResponse()
		s.broadcastC <- oRsp
	}
}

func (s *Sequencer) handleOrderRequests() {
	justReply := s.root && s.leader

	if justReply {
		// in case the sequencer is the root, it will just immediately return with an OrderResponse
		for oReq := range s.oReqCIn {
			sn := s.getAndIncSequenceNum(oReq.NumOfRecords)
			oRsp := &sequencerpb.OrderResponse{
				LocalToken:  oReq.LocalToken,
				GlobalToken: sn,
				Color:       oReq.Color,
				OriginColor: oReq.OriginColor,
			}
			s.broadcastC <- oRsp
		}
		return
	}

	color := s.color

	for oReq := range s.oReqCIn {

		// can just return with responses if leader
		if color == oReq.Color && s.leader {
			// return with orderResponse right away
			sn := s.getAndIncSequenceNum(oReq.NumOfRecords)
			oRsp := &sequencerpb.OrderResponse{
				LocalToken:  oReq.LocalToken,
				GlobalToken: sn,
				Color:       color,
				OriginColor: oReq.OriginColor,
			}
			s.broadcastC <- oRsp
			continue
		}

		s.parentClient.MakeOrderRequest(oReq)
	}
}

func (s *Sequencer) broadcastOrderResponses() {
	for oRsp := range s.broadcastC {
		s.oRspCsMu.RLock()
		for id, oRspC := range s.oRspCs {
			if uint32(id) == oRsp.OriginColor {
				oRspC <- oRsp
			}
		}
		s.oRspCsMu.RUnlock()
	}
}

func (s *Sequencer) forwardOrderResponses(stream sequencerpb.Sequencer_GetOrderServer, oRspC chan *sequencerpb.OrderResponse) {
	for oRsp := range oRspC {
		err := stream.Send(oRsp)
		if err != nil {
			return
		}
	}
}
