package aggregator

import (
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"sync"
	"time"
)

const (
	interval = 100 * time.Microsecond
)

type colorService struct {
	color       uint32
	originColor uint32

	oReqCh chan *sequencerpb.OrderRequest
	oRspCh chan *sequencerpb.OrderResponse

	cache   map[uint64]*sequencerpb.OrderRequest
	cacheMu sync.Mutex
}

func newColorService(color, originColor uint32, onB frame.OnOrderBatch) *colorService {
	cs := new(colorService)
	cs.color = color
	cs.originColor = originColor
	cs.cache = make(map[uint64]*sequencerpb.OrderRequest)
	cs.oReqCh = make(chan *sequencerpb.OrderRequest, 1024)
	cs.oRspCh = make(chan *sequencerpb.OrderResponse, 1024)
	return cs
}

func (cs *colorService) batchOrderRequests(interval time.Duration, onB frame.OnOrderBatch) {
	newBatch := false
	var send <-chan time.Time
	var current *sequencerpb.OrderRequest

	for {
		select {
		case oReq := <-cs.oReqCh:
			if newBatch {
				send = time.After(interval)
				current = oReq
				newBatch = false
				continue
			}
			current.NumOfRecords += 1
		case <-send:
			onB(current)
			newBatch = true
		}
	}
}

func (cs *colorService) MakeOrderRequest(oReq *sequencerpb.OrderRequest) {
	cs.oReqCh <- oReq
}
