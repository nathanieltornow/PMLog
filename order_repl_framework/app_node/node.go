package app_node

import (
	"context"
	"fmt"
	frame "github.com/nathanieltornow/PMLog/order_repl_framework"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/local_node"
	"github.com/nathanieltornow/PMLog/order_repl_framework/app_node/nodepb"
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

	id     uint32
	color  uint32
	ipAddr string

	app frame.Application

	prepareCh     chan *nodepb.Prep
	prepStreams   map[uint32]nodepb.Node_PrepareClient
	prepStreamsMu sync.RWMutex

	waitingORespCh chan *seqpb.OrderResponse

	localNodes   map[uint32]*local_node.LocalNode
	localNodesMu sync.RWMutex

	orderClient *order_client.Client

	ackCh          chan *nodepb.Acknowledgement
	ackBroadCast   map[uint32]chan *nodepb.Acknowledgement
	ackBroadCastMu sync.RWMutex
	numOfPeers     uint32

	helpCtr uint32
}

func NewNode(id, color uint32) (*Node, error) {
	node := new(Node)
	node.id = id
	node.color = color
	node.prepareCh = make(chan *nodepb.Prep, 1024)
	node.prepStreams = make(map[uint32]nodepb.Node_PrepareClient)
	node.waitingORespCh = make(chan *seqpb.OrderResponse, 1024)
	node.localNodes = make(map[uint32]*local_node.LocalNode)
	node.ackCh = make(chan *nodepb.Acknowledgement)
	node.ackBroadCast = make(map[uint32]chan *nodepb.Acknowledgement)
	return node, nil
}

func (n *Node) RegisterApp(app frame.Application) {
	n.app = app
}

func (n *Node) Start(ipAddr string, peerIpAddrs []string, orderIP string) error {
	n.ipAddr = ipAddr
	orderClient, err := order_client.NewClient(orderIP, n.color)
	if err != nil {
		return err
	}
	n.orderClient = orderClient

	errC := make(chan error, 1)
	go func() {
		err := n.startGRPCSever()
		errC <- err
	}()

	go n.broadcastPrepareMsgs()
	go n.handleOrderResponses()

	for _, peerIp := range peerIpAddrs {
		err := n.connectToPeer(peerIp, true)
		if err != nil {
			return err
		}
	}
	n.addNewLocalNode(0)
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

	ackStream, err := client.GetAcknowledgements(context.Background(), &nodepb.AcknowledgeRequest{ID: n.id})
	if err != nil {
		return err
	}
	go n.handleAcks(ackStream)

	if back {
		_, err = client.Register(context.Background(), &nodepb.RegisterRequest{IP: n.ipAddr})
		if err != nil {
			return err
		}
	}
	atomic.AddUint32(&n.numOfPeers, 1)
	n.localNodesMu.RLock()
	for _, lN := range n.localNodes {
		lN.SetPeers(atomic.LoadUint32(&n.numOfPeers))
	}
	n.localNodesMu.RUnlock()
	logrus.Infoln("connected with", peerIP)
	return nil
}

func (n *Node) addNewLocalNode(color uint32) *local_node.LocalNode {
	n.localNodesMu.Lock()
	defer n.localNodesMu.Unlock()
	ln, ok := n.localNodes[color]
	if !ok {
		ln = local_node.NewLocalNode(n.id, atomic.LoadUint32(&n.numOfPeers), n.color, n.app, n.makePrepareMsg, n.orderClient.MakeOrderRequest, n.sendAck)
		n.localNodes[color] = ln
	}
	return ln
}

func (n *Node) getLocalNode(color uint32) *local_node.LocalNode {
	n.localNodesMu.RLock()
	ln, ok := n.localNodes[color]
	n.localNodesMu.RUnlock()
	if !ok {
		return n.addNewLocalNode(color)
	}
	return ln
}
