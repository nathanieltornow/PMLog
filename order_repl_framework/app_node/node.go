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
)

type Node struct {
	nodepb.UnimplementedNodeServer

	app frame.Application

	id     uint32
	color  uint32
	ipAddr string

	ctr uint32

	numOfPeers  uint32
	numOfAcks   map[uint64]uint32
	numOfAcksMu sync.Mutex

	ackChs   map[uint32]chan *nodepb.Ack
	ackChsMu sync.RWMutex

	prepCh        chan *nodepb.Prep
	prepStreams   map[uint32]nodepb.Node_PrepareClient
	prepStreamsMu sync.RWMutex

	comCh        chan *nodepb.Com
	comStreams   map[uint32]nodepb.Node_CommitClient
	comStreamsMu sync.RWMutex

	helpCtr uint32

	orderRespCh <-chan *seqpb.OrderResponse
	orderReqCh  chan<- *seqpb.OrderRequest
}

func NewNode(id, color uint32, app frame.Application) (*Node, error) {
	node := new(Node)
	node.app = app
	node.id = id
	node.color = color
	node.ackChs = make(map[uint32]chan *nodepb.Ack)
	node.prepStreams = make(map[uint32]nodepb.Node_PrepareClient)
	node.prepCh = make(chan *nodepb.Prep, 1024)
	node.numOfAcks = make(map[uint64]uint32)
	node.comCh = make(chan *nodepb.Com)
	node.comStreams = make(map[uint32]nodepb.Node_CommitClient)
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

	errC := make(chan error)
	go func() {
		err := n.startGRPCSever()
		errC <- err
	}()

	go n.broadcastCommitMsgs()
	go n.broadcastPrepareMsgs()
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

func (n *Node) Register(_ context.Context, req *nodepb.RegisterRequest) (*nodepb.Empty, error) {
	if err := n.connectToPeer(req.IP, false); err != nil {
		return nil, err
	}
	return &nodepb.Empty{}, nil
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

	comStream, err := client.Commit(context.Background())
	if err != nil {
		return err
	}
	n.comStreamsMu.Lock()
	n.comStreams[id] = comStream
	n.comStreamsMu.Unlock()

	ackStream, err := client.GetAcks(context.Background(), &nodepb.AckReq{NodeID: n.id})
	if err != nil {
		return err
	}
	go n.receiveAcks(ackStream)

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
