package seq_client

import (
	"context"
	pb "github.com/nathanieltornow/PMLog/sequencer/sequencerpb"
	"google.golang.org/grpc"
)

type Client struct {
	stream pb.Sequencer_GetOrderClient
	oRspC  chan *pb.OrderResponse
	oReqC  chan *pb.OrderRequest
}

func NewClient(IP string) (*Client, error) {
	conn, err := grpc.Dial(IP, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	pbClient := pb.NewSequencerClient(conn)
	stream, err := pbClient.GetOrder(context.Background())
	if err != nil {
		return nil, err
	}
	client := new(Client)
	client.stream = stream
	client.oReqC = make(chan *pb.OrderRequest, 1024)
	client.oRspC = make(chan *pb.OrderResponse, 1024)
	go client.sendOReqs()
	go client.receiveORsps()
	return client, nil
}

func (c *Client) MakeOrderRequests() chan<- *pb.OrderRequest {
	return c.oReqC
}

func (c *Client) GetOrderResponses() <-chan *pb.OrderResponse {
	return c.oRspC
}

func (c *Client) Stop() error {
	err := c.stream.CloseSend()
	if err != nil {
		return err
	}
	close(c.oRspC)
	close(c.oReqC)
	return nil
}

func (c *Client) sendOReqs() {
	for oReq := range c.oReqC {
		err := c.stream.Send(oReq)
		if err != nil {
			return
		}
	}
}

func (c *Client) receiveORsps() {
	for {
		rsp, err := c.stream.Recv()
		if err != nil {
			return
		}
		c.oRspC <- rsp
	}
}
