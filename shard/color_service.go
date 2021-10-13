package shard

import (
	"github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"github.com/nathanieltornow/PMLog/shard/shardpb"
	"github.com/nathanieltornow/PMLog/storage"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type colorService struct {
	mu sync.Mutex

	oRspCh chan *sequencerpb.OrderResponse
	oReqCh chan uint64

	appReqCh chan *shardpb.AppendRequest

	log storage.Log

	waitingAppends map[uint64]chan bool
}

func newColorService(color, originColor uint32, log storage.Log, interval time.Duration,
	outOReqCh chan *sequencerpb.OrderRequest, outAppRspCh chan *shardpb.AppendResponse) *colorService {

	cs := new(colorService)
	cs.oRspCh = make(chan *sequencerpb.OrderResponse, 2048)
	cs.oReqCh = make(chan uint64, 2048)
	cs.appReqCh = make(chan *shardpb.AppendRequest, 2048)
	cs.waitingAppends = make(map[uint64]chan bool)
	cs.log = log

	go cs.handleAppendRequests()
	go cs.handleOrderResponses(outAppRspCh)
	go cs.batchOrderRequests(color, originColor, interval, outOReqCh)
	return cs
}

// ----- API of the color-service

func (cs *colorService) insertAppendRequest(appReq *shardpb.AppendRequest) {
	cs.appReqCh <- appReq
}

func (cs *colorService) insertOrderResponse(oRsp *sequencerpb.OrderResponse) {
	cs.oRspCh <- oRsp
}

func (cs *colorService) read(gsn uint64) (*shardpb.ReadResponse, error) {
	record, err := cs.log.Read(gsn)
	if err != nil {
		return nil, err
	}
	return &shardpb.ReadResponse{Gsn: gsn, Record: record}, nil
}

func (cs *colorService) trim(gsn uint64) (*shardpb.Ok, error) {
	err := cs.log.Trim(gsn)
	if err != nil {
		return nil, err
	}
	return &shardpb.Ok{}, nil
}

// ----- handler threads

func (cs *colorService) handleAppendRequests() {
	for appReq := range cs.appReqCh {
		cs.append(appReq)
		if appReq.Responsible {
			cs.oReqCh <- appReq.Token
		}
	}
}

func (cs *colorService) handleOrderResponses(outAppRspCh chan *shardpb.AppendResponse) {
	for oRsp := range cs.oRspCh {
		for i, token := range oRsp.Tokens {
			cs.checkForAppend(token)
			err := cs.log.Commit(token, oRsp.Gsn+uint64(i))
			if err != nil {
				logrus.Fatalln(err)
			}
			outAppRspCh <- &shardpb.AppendResponse{Token: token, Gsn: oRsp.Gsn + uint64(i)}
		}
	}
}

func (cs *colorService) batchOrderRequests(color, originColor uint32,
	interval time.Duration, outOReqCh chan *sequencerpb.OrderRequest) {

	var send <-chan time.Time

	newBatch := true
	var currentTokens []uint64

	for {
		select {
		case token := <-cs.oReqCh:
			if newBatch {
				currentTokens = make([]uint64, 0)
				send = time.After(interval)
				newBatch = false
			}
			currentTokens = append(currentTokens, token)
		case <-send:
			if newBatch {
				continue
			}
			outOReqCh <- &sequencerpb.OrderRequest{
				OriginColor: originColor,
				Color:       color,
				Tokens:      currentTokens,
			}
			newBatch = true
		}
	}
}

// ---- helper functions for appending

func (cs *colorService) append(appReq *shardpb.AppendRequest) {
	err := cs.log.Append(appReq.Record, appReq.Token)
	if err != nil {
		logrus.Fatalln(err)
	}
	cs.mu.Lock()
	waitCh, ok := cs.waitingAppends[appReq.Token]
	if !ok {
		waitCh = make(chan bool, 1)
		cs.waitingAppends[appReq.Token] = waitCh
	}
	if len(waitCh) == 0 {
		waitCh <- true
	}
	cs.mu.Unlock()
}

func (cs *colorService) checkForAppend(token uint64) {
	cs.mu.Lock()
	waitCh, ok := cs.waitingAppends[token]
	if !ok {
		waitCh = make(chan bool, 1)
		cs.waitingAppends[token] = waitCh
	}
	cs.mu.Unlock()

	defer func() {
		cs.mu.Lock()
		delete(cs.waitingAppends, token)
		cs.mu.Unlock()
	}()
	<-waitCh
}
