// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.21.12
// source: smtp/smtp.proto

package smtp

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	SmtpService_SendEmail_FullMethodName = "/smtp.SmtpService/SendEmail"
)

// SmtpServiceClient is the client API for SmtpService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SmtpServiceClient interface {
	SendEmail(ctx context.Context, in *Email, opts ...grpc.CallOption) (*SendEmailResponse, error)
}

type smtpServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSmtpServiceClient(cc grpc.ClientConnInterface) SmtpServiceClient {
	return &smtpServiceClient{cc}
}

func (c *smtpServiceClient) SendEmail(ctx context.Context, in *Email, opts ...grpc.CallOption) (*SendEmailResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendEmailResponse)
	err := c.cc.Invoke(ctx, SmtpService_SendEmail_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SmtpServiceServer is the server API for SmtpService service.
// All implementations must embed UnimplementedSmtpServiceServer
// for forward compatibility.
type SmtpServiceServer interface {
	SendEmail(context.Context, *Email) (*SendEmailResponse, error)
	mustEmbedUnimplementedSmtpServiceServer()
}

// UnimplementedSmtpServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSmtpServiceServer struct{}

func (UnimplementedSmtpServiceServer) SendEmail(context.Context, *Email) (*SendEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmail not implemented")
}
func (UnimplementedSmtpServiceServer) mustEmbedUnimplementedSmtpServiceServer() {}
func (UnimplementedSmtpServiceServer) testEmbeddedByValue()                     {}

// UnsafeSmtpServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SmtpServiceServer will
// result in compilation errors.
type UnsafeSmtpServiceServer interface {
	mustEmbedUnimplementedSmtpServiceServer()
}

func RegisterSmtpServiceServer(s grpc.ServiceRegistrar, srv SmtpServiceServer) {
	// If the following call pancis, it indicates UnimplementedSmtpServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SmtpService_ServiceDesc, srv)
}

func _SmtpService_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Email)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SmtpServiceServer).SendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SmtpService_SendEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SmtpServiceServer).SendEmail(ctx, req.(*Email))
	}
	return interceptor(ctx, in, info, handler)
}

// SmtpService_ServiceDesc is the grpc.ServiceDesc for SmtpService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SmtpService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smtp.SmtpService",
	HandlerType: (*SmtpServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmail",
			Handler:    _SmtpService_SendEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "smtp/smtp.proto",
}
