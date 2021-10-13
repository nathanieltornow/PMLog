package local_sequencer

//
//import (
//	"github.com/nathanieltornow/PMLog/client_centric/sequencer/sequencerpb"
//	"github.com/sirupsen/logrus"
//	"sync"
//)
//
//type orderRequestCache struct {
//	mu    sync.Mutex
//	ctr   uint64
//	store map[uint64]*sequencerpb.OrderRequest
//}
//
//func newOrderRequestCache() *orderRequestCache {
//	orc := new(orderRequestCache)
//	orc.store = make(map[uint64]*sequencerpb.OrderRequest)
//	return orc
//}
//
//func (orc *orderRequestCache) put(oReq *sequencerpb.OrderRequest) uint64 {
//	orc.mu.Lock()
//	defer orc.mu.Unlock()
//	orc.ctr++
//	orc.store[orc.ctr] = oReq
//	return orc.ctr
//}
//
//func (orc *orderRequestCache) get(oRsp *sequencerpb.OrderResponse) []*sequencerpb.OrderResponse {
//	res := make([]*sequencerpb.OrderResponse, 0)
//	gsn := oRsp.Gsn
//	for i := uint64(0); i < uint64(oRsp.NumOfRecords); i++ {
//		orc.mu.Lock()
//		oReq, ok := orc.store[oRsp.Lsn+i]
//		if !ok {
//			logrus.Fatalln("ups")
//		}
//		delete(orc.store, oRsp.Lsn+i)
//		orc.mu.Unlock()
//
//		newORsp := &sequencerpb.OrderResponse{
//			Lsn:          oReq.Lsn,
//			Color:        oReq.Color,
//			Gsn:          gsn,
//			NumOfRecords: oReq.NumOfRecords,
//			OriginColor:  oReq.OriginColor,
//		}
//		res = append(res, newORsp)
//
//		gsn += uint64(oReq.NumOfRecords)
//	}
//	return res
//}
