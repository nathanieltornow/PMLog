package aggregator

import (
	"context"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
)

type Aggregator struct {
	color       uint32
	originColor uint32
	oRspCh      chan *sequencerpb.OrderResponse

	colorServices map[uint32]*colorService
	mu            sync.RWMutex

	stream sequencerpb.Sequencer_GetOrderClient
}

func NewAggregator(parentIP string, color, originColor uint32) (*Aggregator, error) {
	conn, err := grpc.Dial(parentIP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	pbClient := sequencerpb.NewSequencerClient(conn)
	stream, err := pbClient.GetOrder(context.Background())
	if err != nil {
		return nil, err
	}
	a := new(Aggregator)
	a.color = color
	a.originColor = originColor
	a.stream = stream
	a.oRspCh = make(chan *sequencerpb.OrderResponse, 1024)
	a.colorServices = make(map[uint32]*colorService)
	return a, nil
}

func (a *Aggregator) MakeOrderRequest(oReq *sequencerpb.OrderRequest) {
	a.mu.RLock()
	cs, ok := a.colorServices[oReq.Color]
	a.mu.RUnlock()
	if !ok {
		cs = newColorService(a.color, a.originColor, a.sendOrderRequest)
		a.colorServices[oReq.Color] = cs
	}
}

func (a *Aggregator) GetNextOrderResponse() *sequencerpb.OrderResponse {
	return <-a.oRspCh
}

func (a *Aggregator) receiveOrderResponses() {
	for {
		oRsp, err := a.stream.Recv()
		if err != nil {
			logrus.Fatalln(err)
		}
		a.oRspCh <- oRsp
	}
}

func (a *Aggregator) sendOrderRequest(oReq *sequencerpb.OrderRequest) {
	if err := a.stream.Send(oReq); err != nil {
		logrus.Fatalln(err)
	}
}
