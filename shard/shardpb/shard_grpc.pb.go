// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package shardpb

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

// NodeClient is the client API for Node service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*OK, error)
	Replicate(ctx context.Context, opts ...grpc.CallOption) (Node_ReplicateClient, error)
	Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error)
}

type nodeClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeClient(cc grpc.ClientConnInterface) NodeClient {
	return &nodeClient{cc}
}

func (c *nodeClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*OK, error) {
	out := new(OK)
	err := c.cc.Invoke(ctx, "/sequencer.Node/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nodeClient) Replicate(ctx context.Context, opts ...grpc.CallOption) (Node_ReplicateClient, error) {
	stream, err := c.cc.NewStream(ctx, &Node_ServiceDesc.Streams[0], "/sequencer.Node/Replicate", opts...)
	if err != nil {
		return nil, err
	}
	x := &nodeReplicateClient{stream}
	return x, nil
}

type Node_ReplicateClient interface {
	Send(*ReplicaMessage) error
	Recv() (*ReplicaMessage, error)
	grpc.ClientStream
}

type nodeReplicateClient struct {
	grpc.ClientStream
}

func (x *nodeReplicateClient) Send(m *ReplicaMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *nodeReplicateClient) Recv() (*ReplicaMessage, error) {
	m := new(ReplicaMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *nodeClient) Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error) {
	out := new(AppendResponse)
	err := c.cc.Invoke(ctx, "/sequencer.Node/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeServer is the server API for Node service.
// All implementations must embed UnimplementedNodeServer
// for forward compatibility
type NodeServer interface {
	Register(context.Context, *RegisterRequest) (*OK, error)
	Replicate(Node_ReplicateServer) error
	Append(context.Context, *AppendRequest) (*AppendResponse, error)
	mustEmbedUnimplementedNodeServer()
}

// UnimplementedNodeServer must be embedded to have forward compatible implementations.
type UnimplementedNodeServer struct {
}

func (UnimplementedNodeServer) Register(context.Context, *RegisterRequest) (*OK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedNodeServer) Replicate(Node_ReplicateServer) error {
	return status.Errorf(codes.Unimplemented, "method Replicate not implemented")
}
func (UnimplementedNodeServer) Append(context.Context, *AppendRequest) (*AppendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedNodeServer) mustEmbedUnimplementedNodeServer() {}

// UnsafeNodeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeServer will
// result in compilation errors.
type UnsafeNodeServer interface {
	mustEmbedUnimplementedNodeServer()
}

func RegisterNodeServer(s grpc.ServiceRegistrar, srv NodeServer) {
	s.RegisterService(&Node_ServiceDesc, srv)
}

func _Node_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequencer.Node/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Node_Replicate_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(NodeServer).Replicate(&nodeReplicateServer{stream})
}

type Node_ReplicateServer interface {
	Send(*ReplicaMessage) error
	Recv() (*ReplicaMessage, error)
	grpc.ServerStream
}

type nodeReplicateServer struct {
	grpc.ServerStream
}

func (x *nodeReplicateServer) Send(m *ReplicaMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *nodeReplicateServer) Recv() (*ReplicaMessage, error) {
	m := new(ReplicaMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Node_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sequencer.Node/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeServer).Append(ctx, req.(*AppendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Node_ServiceDesc is the grpc.ServiceDesc for Node service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Node_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sequencer.Node",
	HandlerType: (*NodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Node_Register_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _Node_Append_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Replicate",
			Handler:       _Node_Replicate_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "shardpb/shard.proto",
}