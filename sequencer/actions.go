package sequencer

import (
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
)

func (s *Sequencer) handleOrderResponses() {
	for {
		oRsp := s.parentClient.GetNextOrderResponse()
		oRsps := s.getColorService(oRsp.Color).getOrderResponses(oRsp)
		for _, oRsp := range oRsps {
			s.broadcastCh <- oRsp
		}
	}
}

func (s *Sequencer) forwardOrderRequests() {
	for oReq := range s.oReqCh {
		s.parentClient.MakeOrderRequest(oReq)
	}
}

func (s *Sequencer) broadcastOrderResponses() {
	for oRsp := range s.broadcastCh {
		s.oRspCsMu.RLock()
		for ict, oRspC := range s.oRspCs {
			if ict.color == oRsp.OriginColor {
				oRspC <- oRsp
			}
		}
		s.oRspCsMu.RUnlock()
	}

}

func forwardOrderResponses(stream sequencerpb.Sequencer_GetOrderServer, oRspC chan *sequencerpb.OrderResponse) {
	for oRsp := range oRspC {
		err := stream.Send(oRsp)
		if err != nil {
			return
		}
	}
}
