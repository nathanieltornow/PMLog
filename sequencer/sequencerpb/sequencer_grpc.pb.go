// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package sequencerpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SequencerClient is the client API for Sequencer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SequencerClient interface {
	GetOrder(ctx context.Context, opts ...grpc.CallOption) (Sequencer_GetOrderClient, error)
}

type sequencerClient struct {
	cc grpc.ClientConnInterface
}

func NewSequencerClient(cc grpc.ClientConnInterface) SequencerClient {
	return &sequencerClient{cc}
}

func (c *sequencerClient) GetOrder(ctx context.Context, opts ...grpc.CallOption) (Sequencer_GetOrderClient, error) {
	stream, err := c.cc.NewStream(ctx, &Sequencer_ServiceDesc.Streams[0], "/sequencer.Sequencer/GetOrder", opts...)
	if err != nil {
		return nil, err
	}
	x := &sequencerGetOrderClient{stream}
	return x, nil
}

type Sequencer_GetOrderClient interface {
	Send(*OrderRequest) error
	Recv() (*OrderResponse, error)
	grpc.ClientStream
}

type sequencerGetOrderClient struct {
	grpc.ClientStream
}

func (x *sequencerGetOrderClient) Send(m *OrderRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *sequencerGetOrderClient) Recv() (*OrderResponse, error) {
	m := new(OrderResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SequencerServer is the server API for Sequencer service.
// All implementations must embed UnimplementedSequencerServer
// for forward compatibility
type SequencerServer interface {
	GetOrder(Sequencer_GetOrderServer) error
	mustEmbedUnimplementedSequencerServer()
}

// UnimplementedSequencerServer must be embedded to have forward compatible implementations.
type UnimplementedSequencerServer struct {
}

func (UnimplementedSequencerServer) GetOrder(Sequencer_GetOrderServer) error {
	return status.Errorf(codes.Unimplemented, "method GetOrder not implemented")
}
func (UnimplementedSequencerServer) mustEmbedUnimplementedSequencerServer() {}

// UnsafeSequencerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SequencerServer will
// result in compilation errors.
type UnsafeSequencerServer interface {
	mustEmbedUnimplementedSequencerServer()
}

func RegisterSequencerServer(s grpc.ServiceRegistrar, srv SequencerServer) {
	s.RegisterService(&Sequencer_ServiceDesc, srv)
}

func _Sequencer_GetOrder_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SequencerServer).GetOrder(&sequencerGetOrderServer{stream})
}

type Sequencer_GetOrderServer interface {
	Send(*OrderResponse) error
	Recv() (*OrderRequest, error)
	grpc.ServerStream
}

type sequencerGetOrderServer struct {
	grpc.ServerStream
}

func (x *sequencerGetOrderServer) Send(m *OrderResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *sequencerGetOrderServer) Recv() (*OrderRequest, error) {
	m := new(OrderRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Sequencer_ServiceDesc is the grpc.ServiceDesc for Sequencer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Sequencer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sequencer.Sequencer",
	HandlerType: (*SequencerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetOrder",
			Handler:       _Sequencer_GetOrder_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "sequencerpb/sequencer.proto",
}