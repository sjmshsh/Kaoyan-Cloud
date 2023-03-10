// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: filter_service.proto

package filter_service_v1

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

// FilterServiceClient is the client API for FilterService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FilterServiceClient interface {
	Filter(ctx context.Context, in *ContentMessage, opts ...grpc.CallOption) (*ContentResponse, error)
}

type filterServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFilterServiceClient(cc grpc.ClientConnInterface) FilterServiceClient {
	return &filterServiceClient{cc}
}

func (c *filterServiceClient) Filter(ctx context.Context, in *ContentMessage, opts ...grpc.CallOption) (*ContentResponse, error) {
	out := new(ContentResponse)
	err := c.cc.Invoke(ctx, "/filter.service.v1.FilterService/Filter", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FilterServiceServer is the server API for FilterService service.
// All implementations must embed UnimplementedFilterServiceServer
// for forward compatibility
type FilterServiceServer interface {
	Filter(context.Context, *ContentMessage) (*ContentResponse, error)
	mustEmbedUnimplementedFilterServiceServer()
}

// UnimplementedFilterServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFilterServiceServer struct {
}

func (UnimplementedFilterServiceServer) Filter(context.Context, *ContentMessage) (*ContentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Filter not implemented")
}
func (UnimplementedFilterServiceServer) mustEmbedUnimplementedFilterServiceServer() {}

// UnsafeFilterServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FilterServiceServer will
// result in compilation errors.
type UnsafeFilterServiceServer interface {
	mustEmbedUnimplementedFilterServiceServer()
}

func RegisterFilterServiceServer(s grpc.ServiceRegistrar, srv FilterServiceServer) {
	s.RegisterService(&FilterService_ServiceDesc, srv)
}

func _FilterService_Filter_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ContentMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FilterServiceServer).Filter(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/filter.service.v1.FilterService/Filter",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FilterServiceServer).Filter(ctx, req.(*ContentMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// FilterService_ServiceDesc is the grpc.ServiceDesc for FilterService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FilterService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "filter.service.v1.FilterService",
	HandlerType: (*FilterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Filter",
			Handler:    _FilterService_Filter_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "filter_service.proto",
}
