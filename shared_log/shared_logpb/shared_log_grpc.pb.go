// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package shared_logpb

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

// SharedLogClient is the client API for SharedLog service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SharedLogClient interface {
	Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error)
	Trim(ctx context.Context, in *TrimRequest, opts ...grpc.CallOption) (*TrimResponse, error)
}

type sharedLogClient struct {
	cc grpc.ClientConnInterface
}

func NewSharedLogClient(cc grpc.ClientConnInterface) SharedLogClient {
	return &sharedLogClient{cc}
}

func (c *sharedLogClient) Append(ctx context.Context, in *AppendRequest, opts ...grpc.CallOption) (*AppendResponse, error) {
	out := new(AppendResponse)
	err := c.cc.Invoke(ctx, "/shared_log.SharedLog/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sharedLogClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*ReadResponse, error) {
	out := new(ReadResponse)
	err := c.cc.Invoke(ctx, "/shared_log.SharedLog/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sharedLogClient) Trim(ctx context.Context, in *TrimRequest, opts ...grpc.CallOption) (*TrimResponse, error) {
	out := new(TrimResponse)
	err := c.cc.Invoke(ctx, "/shared_log.SharedLog/Trim", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SharedLogServer is the server API for SharedLog service.
// All implementations must embed UnimplementedSharedLogServer
// for forward compatibility
type SharedLogServer interface {
	Append(context.Context, *AppendRequest) (*AppendResponse, error)
	Read(context.Context, *ReadRequest) (*ReadResponse, error)
	Trim(context.Context, *TrimRequest) (*TrimResponse, error)
	mustEmbedUnimplementedSharedLogServer()
}

// UnimplementedSharedLogServer must be embedded to have forward compatible implementations.
type UnimplementedSharedLogServer struct {
}

func (UnimplementedSharedLogServer) Append(context.Context, *AppendRequest) (*AppendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedSharedLogServer) Read(context.Context, *ReadRequest) (*ReadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedSharedLogServer) Trim(context.Context, *TrimRequest) (*TrimResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Trim not implemented")
}
func (UnimplementedSharedLogServer) mustEmbedUnimplementedSharedLogServer() {}

// UnsafeSharedLogServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SharedLogServer will
// result in compilation errors.
type UnsafeSharedLogServer interface {
	mustEmbedUnimplementedSharedLogServer()
}

func RegisterSharedLogServer(s grpc.ServiceRegistrar, srv SharedLogServer) {
	s.RegisterService(&SharedLog_ServiceDesc, srv)
}

func _SharedLog_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SharedLogServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shared_log.SharedLog/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SharedLogServer).Append(ctx, req.(*AppendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SharedLog_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SharedLogServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shared_log.SharedLog/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SharedLogServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SharedLog_Trim_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TrimRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SharedLogServer).Trim(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/shared_log.SharedLog/Trim",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SharedLogServer).Trim(ctx, req.(*TrimRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SharedLog_ServiceDesc is the grpc.ServiceDesc for SharedLog service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SharedLog_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "shared_log.SharedLog",
	HandlerType: (*SharedLogServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Append",
			Handler:    _SharedLog_Append_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _SharedLog_Read_Handler,
		},
		{
			MethodName: "Trim",
			Handler:    _SharedLog_Trim_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "shared_logpb/shared_log.proto",
}
