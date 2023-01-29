// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: user_service.proto

package user_service_v1

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

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	CheckFileMd5(ctx context.Context, in *CheckFileMd5Request, opts ...grpc.CallOption) (*CheckFileMd5Response, error)
	UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error)
	GetSign(ctx context.Context, in *GetSignRequest, opts ...grpc.CallOption) (*Response, error)
	CheckIn(ctx context.Context, in *CheckSignRequest, opts ...grpc.CallOption) (*Response, error)
	WatchUv(ctx context.Context, in *WatchUvRequest, opts ...grpc.CallOption) (*WatchUvResponse, error)
	Location(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*Response, error)
	FindFriend(ctx context.Context, in *FindFriendRequest, opts ...grpc.CallOption) (*FindFriendResponse, error)
	PostBlog(ctx context.Context, in *PostBlogRequest, opts ...grpc.CallOption) (*Response, error)
	Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (*Response, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	OList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	Comment(ctx context.Context, in *CommentRequest, opts ...grpc.CallOption) (*CommentResponse, error)
	GetFeedList(ctx context.Context, in *GetFeedListRequest, opts ...grpc.CallOption) (*GetFeedListResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) CheckFileMd5(ctx context.Context, in *CheckFileMd5Request, opts ...grpc.CallOption) (*CheckFileMd5Response, error) {
	out := new(CheckFileMd5Response)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/CheckFileMd5", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error) {
	out := new(UploadFileResponse)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/UploadFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetSign(ctx context.Context, in *GetSignRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/GetSign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) CheckIn(ctx context.Context, in *CheckSignRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/CheckIn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) WatchUv(ctx context.Context, in *WatchUvRequest, opts ...grpc.CallOption) (*WatchUvResponse, error) {
	out := new(WatchUvResponse)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/WatchUv", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Location(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/Location", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) FindFriend(ctx context.Context, in *FindFriendRequest, opts ...grpc.CallOption) (*FindFriendResponse, error) {
	out := new(FindFriendResponse)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/FindFriend", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) PostBlog(ctx context.Context, in *PostBlogRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/PostBlog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Watch(ctx context.Context, in *WatchRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/Watch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) OList(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/OList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Comment(ctx context.Context, in *CommentRequest, opts ...grpc.CallOption) (*CommentResponse, error) {
	out := new(CommentResponse)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/Comment", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetFeedList(ctx context.Context, in *GetFeedListRequest, opts ...grpc.CallOption) (*GetFeedListResponse, error) {
	out := new(GetFeedListResponse)
	err := c.cc.Invoke(ctx, "/user.service.v1.UserService/GetFeedList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	CheckFileMd5(context.Context, *CheckFileMd5Request) (*CheckFileMd5Response, error)
	UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error)
	GetSign(context.Context, *GetSignRequest) (*Response, error)
	CheckIn(context.Context, *CheckSignRequest) (*Response, error)
	WatchUv(context.Context, *WatchUvRequest) (*WatchUvResponse, error)
	Location(context.Context, *LocationRequest) (*Response, error)
	FindFriend(context.Context, *FindFriendRequest) (*FindFriendResponse, error)
	PostBlog(context.Context, *PostBlogRequest) (*Response, error)
	Watch(context.Context, *WatchRequest) (*Response, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	OList(context.Context, *ListRequest) (*ListResponse, error)
	Comment(context.Context, *CommentRequest) (*CommentResponse, error)
	GetFeedList(context.Context, *GetFeedListRequest) (*GetFeedListResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) CheckFileMd5(context.Context, *CheckFileMd5Request) (*CheckFileMd5Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckFileMd5 not implemented")
}
func (UnimplementedUserServiceServer) UploadFile(context.Context, *UploadFileRequest) (*UploadFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedUserServiceServer) GetSign(context.Context, *GetSignRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSign not implemented")
}
func (UnimplementedUserServiceServer) CheckIn(context.Context, *CheckSignRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckIn not implemented")
}
func (UnimplementedUserServiceServer) WatchUv(context.Context, *WatchUvRequest) (*WatchUvResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WatchUv not implemented")
}
func (UnimplementedUserServiceServer) Location(context.Context, *LocationRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Location not implemented")
}
func (UnimplementedUserServiceServer) FindFriend(context.Context, *FindFriendRequest) (*FindFriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindFriend not implemented")
}
func (UnimplementedUserServiceServer) PostBlog(context.Context, *PostBlogRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostBlog not implemented")
}
func (UnimplementedUserServiceServer) Watch(context.Context, *WatchRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Watch not implemented")
}
func (UnimplementedUserServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedUserServiceServer) OList(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method OList not implemented")
}
func (UnimplementedUserServiceServer) Comment(context.Context, *CommentRequest) (*CommentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Comment not implemented")
}
func (UnimplementedUserServiceServer) GetFeedList(context.Context, *GetFeedListRequest) (*GetFeedListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFeedList not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_CheckFileMd5_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckFileMd5Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CheckFileMd5(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/CheckFileMd5",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CheckFileMd5(ctx, req.(*CheckFileMd5Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_UploadFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UploadFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/UploadFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UploadFile(ctx, req.(*UploadFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetSign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetSignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetSign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/GetSign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetSign(ctx, req.(*GetSignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_CheckIn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckSignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).CheckIn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/CheckIn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).CheckIn(ctx, req.(*CheckSignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_WatchUv_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WatchUvRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).WatchUv(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/WatchUv",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).WatchUv(ctx, req.(*WatchUvRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Location_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Location(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/Location",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Location(ctx, req.(*LocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_FindFriend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).FindFriend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/FindFriend",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).FindFriend(ctx, req.(*FindFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_PostBlog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostBlogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).PostBlog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/PostBlog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).PostBlog(ctx, req.(*PostBlogRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Watch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Watch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/Watch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Watch(ctx, req.(*WatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_OList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).OList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/OList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).OList(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Comment_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Comment(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/Comment",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Comment(ctx, req.(*CommentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetFeedList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFeedListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetFeedList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/user.service.v1.UserService/GetFeedList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetFeedList(ctx, req.(*GetFeedListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.service.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckFileMd5",
			Handler:    _UserService_CheckFileMd5_Handler,
		},
		{
			MethodName: "UploadFile",
			Handler:    _UserService_UploadFile_Handler,
		},
		{
			MethodName: "GetSign",
			Handler:    _UserService_GetSign_Handler,
		},
		{
			MethodName: "CheckIn",
			Handler:    _UserService_CheckIn_Handler,
		},
		{
			MethodName: "WatchUv",
			Handler:    _UserService_WatchUv_Handler,
		},
		{
			MethodName: "Location",
			Handler:    _UserService_Location_Handler,
		},
		{
			MethodName: "FindFriend",
			Handler:    _UserService_FindFriend_Handler,
		},
		{
			MethodName: "PostBlog",
			Handler:    _UserService_PostBlog_Handler,
		},
		{
			MethodName: "Watch",
			Handler:    _UserService_Watch_Handler,
		},
		{
			MethodName: "List",
			Handler:    _UserService_List_Handler,
		},
		{
			MethodName: "OList",
			Handler:    _UserService_OList_Handler,
		},
		{
			MethodName: "Comment",
			Handler:    _UserService_Comment_Handler,
		},
		{
			MethodName: "GetFeedList",
			Handler:    _UserService_GetFeedList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user_service.proto",
}