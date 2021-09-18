package app_node

import (
	"context"
	"fmt"
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	order_client "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/client"
	seqpb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	maxMsgSize    = 128
	batchInterval = 10 * time.Microsecond
)

type Node struct {
	nodepb.UnimplementedNodeServer

	app frame.Application

	id      uint32
	color   uint32
	ipAddr  string
	ctr     uint32
	helpCtr uint32

	maxMsgSize int

	batchInterval time.Duration

	prepStreams   map[uint32]nodepb.Node_PrepareClient
	prepStreamsMu sync.RWMutex

	waitingORespCh chan *seqpb.OrderResponse

	orderRespCh <-chan *seqpb.OrderResponse
	orderReqCh  chan<- *seqpb.OrderRequest

	prepMan *prepManager

	prepB *prepBatch
}

func NewNode(id, color uint32, app frame.Application) (*Node, error) {
	node := new(Node)
	node.app = app
	node.id = id
	node.color = color
	node.prepStreams = make(map[uint32]nodepb.Node_PrepareClient)
	node.waitingORespCh = make(chan *seqpb.OrderResponse, 1024)
	node.maxMsgSize = maxMsgSize
	node.batchInterval = batchInterval
	return node, nil
}

func (n *Node) Start(ipAddr string, peerIpAddrs []string, orderIP string) error {
	n.ipAddr = ipAddr
	orderClient, err := order_client.NewClient(orderIP)
	if err != nil {
		return err
	}
	n.orderReqCh = orderClient.MakeOrderRequests()
	n.orderRespCh = orderClient.GetOrderResponses()

	errC := make(chan error, 1)
	go func() {
		err := n.startGRPCSever()
		errC <- err
	}()

	n.prepMan = newPrepManager(n.app)
	n.prepB = newPrepBatch(n.broadcastPrepareMsgs, n.maxMsgSize, n.batchInterval)

	go n.handleAppCommitRequests()
	go n.handleOrderResponses()

	for _, peerIp := range peerIpAddrs {
		err := n.connectToPeer(peerIp, true)
		if err != nil {
			return err
		}
	}
	err = <-errC
	return err
}

func (n *Node) Register(_ context.Context, req *nodepb.RegisterRequest) (*nodepb.Empty, error) {
	if err := n.connectToPeer(req.IP, false); err != nil {
		return nil, err
	}
	return &nodepb.Empty{}, nil
}

func (n *Node) startGRPCSever() error {
	lis, err := net.Listen("tcp", n.ipAddr)
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	nodepb.RegisterNodeServer(server, n)

	logrus.Infoln("Starting node")
	if err := server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start node: %v", err)
	}
	return nil
}

func (n *Node) connectToPeer(peerIP string, back bool) error {
	conn, err := grpc.Dial(peerIP, grpc.WithInsecure())
	if err != nil {
		return err
	}
	client := nodepb.NewNodeClient(conn)

	id := atomic.AddUint32(&n.helpCtr, 1)
	prepStream, err := client.Prepare(context.Background())
	if err != nil {
		return err
	}
	n.prepStreamsMu.Lock()
	n.prepStreams[id] = prepStream
	n.prepStreamsMu.Unlock()

	if back {
		_, err = client.Register(context.Background(), &nodepb.RegisterRequest{IP: n.ipAddr})
		if err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) getNewLocalToken() uint64 {
	ctr := atomic.AddUint32(&n.ctr, 1)
	return (uint64(ctr) << 32) + uint64(n.id)
}
