package client

import (
	"context"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"sync"
	"time"
)

const (
	orderReqInterval = time.Microsecond * 100
)

type Client struct {
	stream sequencerpb.Sequencer_GetOrderClient

	interval time.Duration

	color uint32

	mu    sync.RWMutex
	cache map[uint64][]*sequencerpb.OrderRequest

	oRspC chan *sequencerpb.OrderResponse
	oReqC chan *sequencerpb.OrderRequest

	sendC chan *sequencerpb.BatchedOrderRequest
}

func NewClient(IP string, color uint32, options ...Option) (*Client, error) {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	pbClient := sequencerpb.NewSequencerClient(conn)
	stream, err := pbClient.GetOrder(context.Background())
	if err != nil {
		return nil, err
	}
	client := new(Client)
	client.stream = stream
	client.oReqC = make(chan *sequencerpb.OrderRequest, 1024)
	client.oRspC = make(chan *sequencerpb.OrderResponse, 1024)
	client.sendC = make(chan *sequencerpb.BatchedOrderRequest, 1024)
	client.interval = orderReqInterval
	client.color = color
	client.cache = make(map[uint64][]*sequencerpb.OrderRequest)
	for _, o := range options {
		o(client)
	}
	go client.batchOrderRequests()
	go client.receiveORsps()
	go client.sendBatchedOReqs()
	return client, nil
}

func (c *Client) MakeOrderRequest(oReq *sequencerpb.OrderRequest) {
	c.oReqC <- oReq
}

func (c *Client) GetNextOrderResponse() *sequencerpb.OrderResponse {
	oRsp := <-c.oRspC
	return oRsp
}

func (c *Client) batchOrderRequests() {

	var colorToOReq map[uint32]*sequencerpb.OrderRequest

	newBatch := true
	var send <-chan time.Time

	token := uint32(0)

	for {
		select {
		case oReq := <-c.oReqC:
			if newBatch {
				colorToOReq = make(map[uint32]*sequencerpb.OrderRequest)
				send = time.After(c.interval)
				newBatch = false
			}

			localToken := c.addToCache(token, oReq)

			curOReq, ok := colorToOReq[oReq.Color]
			if !ok {
				colorToOReq[oReq.Color] = &sequencerpb.OrderRequest{
					LocalToken:   localToken,
					NumOfRecords: oReq.NumOfRecords,
					Color:        oReq.Color,
					OriginColor:  c.color,
				}
				continue
			}
			curOReq.NumOfRecords += oReq.NumOfRecords

		case <-send:
			token++
			c.sendC <- mapToBatchedOReq(colorToOReq)
			newBatch = true
		}
	}
}

func (c *Client) sendBatchedOReqs() {
	for bOReq := range c.sendC {
		if err := c.stream.Send(bOReq); err != nil {
			logrus.Fatalln(err)
		}
	}
}

func (c *Client) receiveORsps() {
	for {
		rsp, err := c.stream.Recv()
		if err != nil {
			logrus.Fatalln(err)
		}
		if rsp.OriginColor != c.color {
			continue
		}
		c.enqueueOrderResponses(rsp)
	}
}

func (c *Client) addToCache(token uint32, oReq *sequencerpb.OrderRequest) uint64 {
	colorToken := (uint64(token) << 32) + uint64(oReq.Color)
	c.mu.Lock()
	defer c.mu.Unlock()
	oReqList, ok := c.cache[colorToken]
	if !ok {
		newList := []*sequencerpb.OrderRequest{oReq}
		c.cache[colorToken] = newList
		return colorToken
	}
	oReqList = append(oReqList, oReq)
	return colorToken
}

func (c *Client) enqueueOrderResponses(inORsp *sequencerpb.OrderResponse) {
	oReqList, ok := c.cache[inORsp.LocalToken]
	if !ok {
		logrus.Fatalln("failed to find orderRequests")
	}
	globalToken := inORsp.GlobalToken
	for _, oReq := range oReqList {
		c.oRspC <- &sequencerpb.OrderResponse{
			LocalToken:   oReq.LocalToken,
			NumOfRecords: oReq.NumOfRecords,
			Color:        oReq.Color,
			OriginColor:  oReq.OriginColor,
			GlobalToken:  globalToken,
		}
		globalToken += uint64(oReq.NumOfRecords)
	}
}

func mapToBatchedOReq(oReqMap map[uint32]*sequencerpb.OrderRequest) *sequencerpb.BatchedOrderRequest {
	oReqs := make([]*sequencerpb.OrderRequest, 0)
	for _, oReq := range oReqMap {
		oReqs = append(oReqs, oReq)
	}
	return &sequencerpb.BatchedOrderRequest{OReqs: oReqs}
}
