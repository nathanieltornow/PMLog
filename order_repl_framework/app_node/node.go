package app_node

import (
	"context"
	"fmt"
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/color_service"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
	"github.com/nathanieltornow/PMLog/order_repl_framework/sequencer"
	order_client "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/client"
	seqpb "github.com/nathanieltornow/PMLog/order_repl_framework/sequencer/sequencerpb"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"sync"
	"sync/atomic"
)

type Node struct {
	nodepb.UnimplementedNodeServer

	app frame.Application

	id      uint32
	color   uint32
	ipAddr  string
	helpCtr uint32

	colorServices   map[uint32]*color_service.ColorService
	colorServicesMu sync.RWMutex

	maxMsgSize int

	prepStreams   map[uint32]nodepb.Node_PrepareClient
	prepStreamsMu sync.RWMutex

	waitingORespCh chan *seqpb.OrderResponse

	localSequencer *sequencer.Sequencer
}

func NewNode(id, color uint32, app frame.Application, options ...NodeOption) (*Node, error) {
	node := new(Node)
	node.app = app
	node.id = id
	node.color = color
	node.prepStreams = make(map[uint32]nodepb.Node_PrepareClient)
	node.waitingORespCh = make(chan *seqpb.OrderResponse, 1024)
	node.colorServices = make(map[uint32]*color_service.ColorService)
	for _, o := range options {
		o(node)
	}
	return node, nil
}

func (n *Node) Start(ipAddr string, peerIpAddrs []string, orderIP string) error {
	n.ipAddr = ipAddr
	orderClient, err := order_client.NewClient(orderIP, n.color)
	if err != nil {
		return err
	}

	errC := make(chan error, 1)
	go func() {
		err := n.startGRPCSever()
		errC <- err
	}()

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
