package local_sequencer

//
//import (
//	frame "github.com/nathanieltornow/PMLog/client_centric"
//	"github.com/nathanieltornow/PMLog/client_centric/sequencer/sequencerpb"
//	"time"
//)
//
//const (
//	batchInterval = 4 * time.Second
//)
//
//type LocalSequencer struct {
//	cache  *orderRequestCache
//	oReqCh chan *sequencerpb.OrderRequest
//}
//
//func NewLocalSequencer(onBatch frame.OnOrderFunc) *LocalSequencer {
//	ls := new(LocalSequencer)
//	ls.oReqCh = make(chan *sequencerpb.OrderRequest, 1024)
//	ls.cache = newOrderRequestCache()
//
//	go ls.batchOrderRequests(batchInterval, onBatch)
//	return ls
//}
//
//func (ls *LocalSequencer) MakeOrderRequest(oReq *sequencerpb.OrderRequest) {
//	ls.oReqCh <- oReq
//}
//
//func (ls *LocalSequencer) GetOrderResponses(oRsp *sequencerpb.OrderResponse) []*sequencerpb.OrderResponse {
//	return ls.cache.get(oRsp)
//}
//
//func (ls *LocalSequencer) batchOrderRequests(interval time.Duration, onBatch frame.OnOrderFunc) {
//	newBatch := true
//	var send <-chan time.Time
//	var current *sequencerpb.OrderRequest
//
//	for {
//		select {
//		case oReq := <-ls.oReqCh:
//			lsn := ls.cache.put(oReq)
//			if newBatch {
//				send = time.After(interval)
//				current = &sequencerpb.OrderRequest{Lsn: lsn, Color: oReq.Color, OriginColor: oReq.OriginColor, NumOfRecords: oReq.NumOfRecords}
//				newBatch = false
//				continue
//			}
//			current.NumOfRecords += oReq.NumOfRecords
//		case <-send:
//			onBatch(current)
//			newBatch = true
//		}
//	}
//}
