package sequencer

import (
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
)

func (s *Sequencer) handleOrderResponses() {
	//for {
	//	oRsp := s.parentClient.GetNextOrderResponse()
	//	oRsps := s.getLocalSequencer(oRsp.Color).GetOrderResponses(oRsp)
	//	for _, oRsp := range oRsps {
	//		s.broadcastCh <- oRsp
	//	}
	//}
}

func (s *Sequencer) handleOrderRequests() {
	justReply := s.root

	if justReply {
		// in case the sequencer is the root, it will just immediately return with an OrderResponse
		for oReq := range s.oReqCIn {
			sn := s.getAndIncSequenceNum(uint32(len(oReq.Tokens)))
			oRsp := &sequencerpb.OrderResponse{
				Tokens:      oReq.Tokens,
				Gsn:         sn,
				Color:       oReq.Color,
				OriginColor: oReq.OriginColor,
			}
			s.broadcastCh <- oRsp
		}
		return
	}

	//for oReq := range s.oReqCIn {
	//	s.getLocalSequencer(oReq.Color).MakeOrderRequest(oReq)
	//}

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

func (s *Sequencer) forwardOrderResponses(stream sequencerpb.Sequencer_GetOrderServer, oRspC chan *sequencerpb.OrderResponse) {
	for oRsp := range oRspC {
		err := stream.Send(oRsp)
		if err != nil {
			return
		}
	}
}

//func (s *Sequencer) addNewLocalSequencer(color uint32) *local_sequencer.LocalSequencer {
//	s.localSeqsMu.Lock()
//	defer s.localSeqsMu.Unlock()
//	if localSeq, ok := s.localSeqs[color]; ok {
//		return localSeq
//	}
//	locSeq := local_sequencer.NewLocalSequencer(s.parentClient.MakeOrderRequest)
//	s.localSeqs[color] = locSeq
//	return locSeq
//}

//func (s *Sequencer) getLocalSequencer(color uint32) *local_sequencer.LocalSequencer {
//	s.localSeqsMu.RLock()
//	oReqCache, ok := s.localSeqs[color]
//	s.localSeqsMu.RUnlock()
//	if ok {
//		return oReqCache
//	}
//	return s.addNewLocalSequencer(color)
//}
