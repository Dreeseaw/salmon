// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

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

// RouterServiceClient is the client API for RouterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RouterServiceClient interface {
	// Connect to router
	Connect(ctx context.Context, in *ClientID, opts ...grpc.CallOption) (*SuccessResponse, error)
	// Unconnect from router
	// Client must send own id for success
	Disconnect(ctx context.Context, in *ClientID, opts ...grpc.CallOption) (*SuccessResponse, error)
	// Insert a new object into system
	SendInsert(ctx context.Context, in *InsertCommand, opts ...grpc.CallOption) (*SuccessResponse, error)
	// General query rpc
	// load balance per-required-partition
	SendSelect(ctx context.Context, in *SelectCommand, opts ...grpc.CallOption) (RouterService_SendSelectClient, error)
	// Stream replicated objects to clients for storage
	ReceiveReplicas(ctx context.Context, opts ...grpc.CallOption) (RouterService_ReceiveReplicasClient, error)
	// Stream load-balanced partial commands to clients
	// for processing, read in PartialResults
	ProcessPartials(ctx context.Context, opts ...grpc.CallOption) (RouterService_ProcessPartialsClient, error)
}

type routerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRouterServiceClient(cc grpc.ClientConnInterface) RouterServiceClient {
	return &routerServiceClient{cc}
}

func (c *routerServiceClient) Connect(ctx context.Context, in *ClientID, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := c.cc.Invoke(ctx, "/RouterService.RouterService/Connect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerServiceClient) Disconnect(ctx context.Context, in *ClientID, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := c.cc.Invoke(ctx, "/RouterService.RouterService/Disconnect", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerServiceClient) SendInsert(ctx context.Context, in *InsertCommand, opts ...grpc.CallOption) (*SuccessResponse, error) {
	out := new(SuccessResponse)
	err := c.cc.Invoke(ctx, "/RouterService.RouterService/SendInsert", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *routerServiceClient) SendSelect(ctx context.Context, in *SelectCommand, opts ...grpc.CallOption) (RouterService_SendSelectClient, error) {
	stream, err := c.cc.NewStream(ctx, &RouterService_ServiceDesc.Streams[0], "/RouterService.RouterService/SendSelect", opts...)
	if err != nil {
		return nil, err
	}
	x := &routerServiceSendSelectClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RouterService_SendSelectClient interface {
	Recv() (*PartialResult, error)
	grpc.ClientStream
}

type routerServiceSendSelectClient struct {
	grpc.ClientStream
}

func (x *routerServiceSendSelectClient) Recv() (*PartialResult, error) {
	m := new(PartialResult)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *routerServiceClient) ReceiveReplicas(ctx context.Context, opts ...grpc.CallOption) (RouterService_ReceiveReplicasClient, error) {
	stream, err := c.cc.NewStream(ctx, &RouterService_ServiceDesc.Streams[1], "/RouterService.RouterService/ReceiveReplicas", opts...)
	if err != nil {
		return nil, err
	}
	x := &routerServiceReceiveReplicasClient{stream}
	return x, nil
}

type RouterService_ReceiveReplicasClient interface {
	Send(*SuccessResponse) error
	Recv() (*InsertCommand, error)
	grpc.ClientStream
}

type routerServiceReceiveReplicasClient struct {
	grpc.ClientStream
}

func (x *routerServiceReceiveReplicasClient) Send(m *SuccessResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *routerServiceReceiveReplicasClient) Recv() (*InsertCommand, error) {
	m := new(InsertCommand)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *routerServiceClient) ProcessPartials(ctx context.Context, opts ...grpc.CallOption) (RouterService_ProcessPartialsClient, error) {
	stream, err := c.cc.NewStream(ctx, &RouterService_ServiceDesc.Streams[2], "/RouterService.RouterService/ProcessPartials", opts...)
	if err != nil {
		return nil, err
	}
	x := &routerServiceProcessPartialsClient{stream}
	return x, nil
}

type RouterService_ProcessPartialsClient interface {
	Send(*PartialResult) error
	Recv() (*SelectCommand, error)
	grpc.ClientStream
}

type routerServiceProcessPartialsClient struct {
	grpc.ClientStream
}

func (x *routerServiceProcessPartialsClient) Send(m *PartialResult) error {
	return x.ClientStream.SendMsg(m)
}

func (x *routerServiceProcessPartialsClient) Recv() (*SelectCommand, error) {
	m := new(SelectCommand)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouterServiceServer is the server API for RouterService service.
// All implementations must embed UnimplementedRouterServiceServer
// for forward compatibility
type RouterServiceServer interface {
	// Connect to router
	Connect(context.Context, *ClientID) (*SuccessResponse, error)
	// Unconnect from router
	// Client must send own id for success
	Disconnect(context.Context, *ClientID) (*SuccessResponse, error)
	// Insert a new object into system
	SendInsert(context.Context, *InsertCommand) (*SuccessResponse, error)
	// General query rpc
	// load balance per-required-partition
	SendSelect(*SelectCommand, RouterService_SendSelectServer) error
	// Stream replicated objects to clients for storage
	ReceiveReplicas(RouterService_ReceiveReplicasServer) error
	// Stream load-balanced partial commands to clients
	// for processing, read in PartialResults
	ProcessPartials(RouterService_ProcessPartialsServer) error
	mustEmbedUnimplementedRouterServiceServer()
}

// UnimplementedRouterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRouterServiceServer struct {
}

func (UnimplementedRouterServiceServer) Connect(context.Context, *ClientID) (*SuccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedRouterServiceServer) Disconnect(context.Context, *ClientID) (*SuccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disconnect not implemented")
}
func (UnimplementedRouterServiceServer) SendInsert(context.Context, *InsertCommand) (*SuccessResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendInsert not implemented")
}
func (UnimplementedRouterServiceServer) SendSelect(*SelectCommand, RouterService_SendSelectServer) error {
	return status.Errorf(codes.Unimplemented, "method SendSelect not implemented")
}
func (UnimplementedRouterServiceServer) ReceiveReplicas(RouterService_ReceiveReplicasServer) error {
	return status.Errorf(codes.Unimplemented, "method ReceiveReplicas not implemented")
}
func (UnimplementedRouterServiceServer) ProcessPartials(RouterService_ProcessPartialsServer) error {
	return status.Errorf(codes.Unimplemented, "method ProcessPartials not implemented")
}
func (UnimplementedRouterServiceServer) mustEmbedUnimplementedRouterServiceServer() {}

// UnsafeRouterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RouterServiceServer will
// result in compilation errors.
type UnsafeRouterServiceServer interface {
	mustEmbedUnimplementedRouterServiceServer()
}

func RegisterRouterServiceServer(s grpc.ServiceRegistrar, srv RouterServiceServer) {
	s.RegisterService(&RouterService_ServiceDesc, srv)
}

func _RouterService_Connect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouterServiceServer).Connect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RouterService.RouterService/Connect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouterServiceServer).Connect(ctx, req.(*ClientID))
	}
	return interceptor(ctx, in, info, handler)
}

func _RouterService_Disconnect_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouterServiceServer).Disconnect(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RouterService.RouterService/Disconnect",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouterServiceServer).Disconnect(ctx, req.(*ClientID))
	}
	return interceptor(ctx, in, info, handler)
}

func _RouterService_SendInsert_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InsertCommand)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RouterServiceServer).SendInsert(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RouterService.RouterService/SendInsert",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RouterServiceServer).SendInsert(ctx, req.(*InsertCommand))
	}
	return interceptor(ctx, in, info, handler)
}

func _RouterService_SendSelect_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SelectCommand)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RouterServiceServer).SendSelect(m, &routerServiceSendSelectServer{stream})
}

type RouterService_SendSelectServer interface {
	Send(*PartialResult) error
	grpc.ServerStream
}

type routerServiceSendSelectServer struct {
	grpc.ServerStream
}

func (x *routerServiceSendSelectServer) Send(m *PartialResult) error {
	return x.ServerStream.SendMsg(m)
}

func _RouterService_ReceiveReplicas_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RouterServiceServer).ReceiveReplicas(&routerServiceReceiveReplicasServer{stream})
}

type RouterService_ReceiveReplicasServer interface {
	Send(*InsertCommand) error
	Recv() (*SuccessResponse, error)
	grpc.ServerStream
}

type routerServiceReceiveReplicasServer struct {
	grpc.ServerStream
}

func (x *routerServiceReceiveReplicasServer) Send(m *InsertCommand) error {
	return x.ServerStream.SendMsg(m)
}

func (x *routerServiceReceiveReplicasServer) Recv() (*SuccessResponse, error) {
	m := new(SuccessResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _RouterService_ProcessPartials_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RouterServiceServer).ProcessPartials(&routerServiceProcessPartialsServer{stream})
}

type RouterService_ProcessPartialsServer interface {
	Send(*SelectCommand) error
	Recv() (*PartialResult, error)
	grpc.ServerStream
}

type routerServiceProcessPartialsServer struct {
	grpc.ServerStream
}

func (x *routerServiceProcessPartialsServer) Send(m *SelectCommand) error {
	return x.ServerStream.SendMsg(m)
}

func (x *routerServiceProcessPartialsServer) Recv() (*PartialResult, error) {
	m := new(PartialResult)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// RouterService_ServiceDesc is the grpc.ServiceDesc for RouterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RouterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RouterService.RouterService",
	HandlerType: (*RouterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Connect",
			Handler:    _RouterService_Connect_Handler,
		},
		{
			MethodName: "Disconnect",
			Handler:    _RouterService_Disconnect_Handler,
		},
		{
			MethodName: "SendInsert",
			Handler:    _RouterService_SendInsert_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendSelect",
			Handler:       _RouterService_SendSelect_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ReceiveReplicas",
			Handler:       _RouterService_ReceiveReplicas_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "ProcessPartials",
			Handler:       _RouterService_ProcessPartials_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "grpc/router_service.proto",
}
