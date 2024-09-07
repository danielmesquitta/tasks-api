// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: health_check_service.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	HealthCheckService_Check_FullMethodName = "/tasksapi.HealthCheckService/Check"
)

// HealthCheckServiceClient is the client API for HealthCheckService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HealthCheckServiceClient interface {
	Check(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type healthCheckServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewHealthCheckServiceClient(cc grpc.ClientConnInterface) HealthCheckServiceClient {
	return &healthCheckServiceClient{cc}
}

func (c *healthCheckServiceClient) Check(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, HealthCheckService_Check_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HealthCheckServiceServer is the server API for HealthCheckService service.
// All implementations must embed UnimplementedHealthCheckServiceServer
// for forward compatibility.
type HealthCheckServiceServer interface {
	Check(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	mustEmbedUnimplementedHealthCheckServiceServer()
}

// UnimplementedHealthCheckServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedHealthCheckServiceServer struct{}

func (UnimplementedHealthCheckServiceServer) Check(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedHealthCheckServiceServer) mustEmbedUnimplementedHealthCheckServiceServer() {}
func (UnimplementedHealthCheckServiceServer) testEmbeddedByValue()                            {}

// UnsafeHealthCheckServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HealthCheckServiceServer will
// result in compilation errors.
type UnsafeHealthCheckServiceServer interface {
	mustEmbedUnimplementedHealthCheckServiceServer()
}

func RegisterHealthCheckServiceServer(s grpc.ServiceRegistrar, srv HealthCheckServiceServer) {
	// If the following call pancis, it indicates UnimplementedHealthCheckServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&HealthCheckService_ServiceDesc, srv)
}

func _HealthCheckService_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HealthCheckServiceServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: HealthCheckService_Check_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HealthCheckServiceServer).Check(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// HealthCheckService_ServiceDesc is the grpc.ServiceDesc for HealthCheckService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var HealthCheckService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "tasksapi.HealthCheckService",
	HandlerType: (*HealthCheckServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Check",
			Handler:    _HealthCheckService_Check_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "health_check_service.proto",
}