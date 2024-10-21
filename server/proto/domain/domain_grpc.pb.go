// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: domain/domain.proto

package domain

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

// DomainServiceClient is the client API for DomainService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DomainServiceClient interface {
	QueueDNSVerification(ctx context.Context, in *QueueDNSVerificationRequest, opts ...grpc.CallOption) (*QueueDNSVerificationResponse, error)
}

type domainServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDomainServiceClient(cc grpc.ClientConnInterface) DomainServiceClient {
	return &domainServiceClient{cc}
}

func (c *domainServiceClient) QueueDNSVerification(ctx context.Context, in *QueueDNSVerificationRequest, opts ...grpc.CallOption) (*QueueDNSVerificationResponse, error) {
	out := new(QueueDNSVerificationResponse)
	err := c.cc.Invoke(ctx, "/domain.DomainService/QueueDNSVerification", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DomainServiceServer is the server API for DomainService service.
// All implementations must embed UnimplementedDomainServiceServer
// for forward compatibility
type DomainServiceServer interface {
	QueueDNSVerification(context.Context, *QueueDNSVerificationRequest) (*QueueDNSVerificationResponse, error)
	mustEmbedUnimplementedDomainServiceServer()
}

// UnimplementedDomainServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDomainServiceServer struct {
}

func (UnimplementedDomainServiceServer) QueueDNSVerification(context.Context, *QueueDNSVerificationRequest) (*QueueDNSVerificationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueueDNSVerification not implemented")
}
func (UnimplementedDomainServiceServer) mustEmbedUnimplementedDomainServiceServer() {}

// UnsafeDomainServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DomainServiceServer will
// result in compilation errors.
type UnsafeDomainServiceServer interface {
	mustEmbedUnimplementedDomainServiceServer()
}

func RegisterDomainServiceServer(s grpc.ServiceRegistrar, srv DomainServiceServer) {
	s.RegisterService(&DomainService_ServiceDesc, srv)
}

func _DomainService_QueueDNSVerification_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueueDNSVerificationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DomainServiceServer).QueueDNSVerification(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/domain.DomainService/QueueDNSVerification",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DomainServiceServer).QueueDNSVerification(ctx, req.(*QueueDNSVerificationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DomainService_ServiceDesc is the grpc.ServiceDesc for DomainService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DomainService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "domain.DomainService",
	HandlerType: (*DomainServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueueDNSVerification",
			Handler:    _DomainService_QueueDNSVerification_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "domain/domain.proto",
}
