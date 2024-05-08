// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: room.proto

package gen

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

const (
	RoomService_CreateRoom_FullMethodName       = "/watch2gather.proto.roomapi.RoomService/CreateRoom"
	RoomService_GetRoomsByUser_FullMethodName   = "/watch2gather.proto.roomapi.RoomService/GetRoomsByUser"
	RoomService_GetUserIDsByRoom_FullMethodName = "/watch2gather.proto.roomapi.RoomService/GetUserIDsByRoom"
	RoomService_InviteToRoom_FullMethodName     = "/watch2gather.proto.roomapi.RoomService/InviteToRoom"
	RoomService_EnterRoom_FullMethodName        = "/watch2gather.proto.roomapi.RoomService/EnterRoom"
	RoomService_SendMessage_FullMethodName      = "/watch2gather.proto.roomapi.RoomService/SendMessage"
	RoomService_UpdateRoom_FullMethodName       = "/watch2gather.proto.roomapi.RoomService/UpdateRoom"
	RoomService_DeleteRoom_FullMethodName       = "/watch2gather.proto.roomapi.RoomService/DeleteRoom"
)

// RoomServiceClient is the client API for RoomService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RoomServiceClient interface {
	CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error)
	GetRoomsByUser(ctx context.Context, in *GetRoomsByUserRequest, opts ...grpc.CallOption) (*GetRoomsByUserResponse, error)
	GetUserIDsByRoom(ctx context.Context, in *GetUserIDsByRoomRequest, opts ...grpc.CallOption) (*GetUserIDsByRoomResponse, error)
	InviteToRoom(ctx context.Context, in *InviteToRoomRequest, opts ...grpc.CallOption) (*InviteToRoomResponse, error)
	EnterRoom(ctx context.Context, in *EnterRoomRequest, opts ...grpc.CallOption) (RoomService_EnterRoomClient, error)
	SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*SendMessageResponse, error)
	UpdateRoom(ctx context.Context, in *UpdateRoomRequest, opts ...grpc.CallOption) (*UpdateRoomResponse, error)
	DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*DeleteRoomResponse, error)
}

type roomServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRoomServiceClient(cc grpc.ClientConnInterface) RoomServiceClient {
	return &roomServiceClient{cc}
}

func (c *roomServiceClient) CreateRoom(ctx context.Context, in *CreateRoomRequest, opts ...grpc.CallOption) (*CreateRoomResponse, error) {
	out := new(CreateRoomResponse)
	err := c.cc.Invoke(ctx, RoomService_CreateRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) GetRoomsByUser(ctx context.Context, in *GetRoomsByUserRequest, opts ...grpc.CallOption) (*GetRoomsByUserResponse, error) {
	out := new(GetRoomsByUserResponse)
	err := c.cc.Invoke(ctx, RoomService_GetRoomsByUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) GetUserIDsByRoom(ctx context.Context, in *GetUserIDsByRoomRequest, opts ...grpc.CallOption) (*GetUserIDsByRoomResponse, error) {
	out := new(GetUserIDsByRoomResponse)
	err := c.cc.Invoke(ctx, RoomService_GetUserIDsByRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) InviteToRoom(ctx context.Context, in *InviteToRoomRequest, opts ...grpc.CallOption) (*InviteToRoomResponse, error) {
	out := new(InviteToRoomResponse)
	err := c.cc.Invoke(ctx, RoomService_InviteToRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) EnterRoom(ctx context.Context, in *EnterRoomRequest, opts ...grpc.CallOption) (RoomService_EnterRoomClient, error) {
	stream, err := c.cc.NewStream(ctx, &RoomService_ServiceDesc.Streams[0], RoomService_EnterRoom_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &roomServiceEnterRoomClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type RoomService_EnterRoomClient interface {
	Recv() (*Message, error)
	grpc.ClientStream
}

type roomServiceEnterRoomClient struct {
	grpc.ClientStream
}

func (x *roomServiceEnterRoomClient) Recv() (*Message, error) {
	m := new(Message)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *roomServiceClient) SendMessage(ctx context.Context, in *Message, opts ...grpc.CallOption) (*SendMessageResponse, error) {
	out := new(SendMessageResponse)
	err := c.cc.Invoke(ctx, RoomService_SendMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) UpdateRoom(ctx context.Context, in *UpdateRoomRequest, opts ...grpc.CallOption) (*UpdateRoomResponse, error) {
	out := new(UpdateRoomResponse)
	err := c.cc.Invoke(ctx, RoomService_UpdateRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *roomServiceClient) DeleteRoom(ctx context.Context, in *DeleteRoomRequest, opts ...grpc.CallOption) (*DeleteRoomResponse, error) {
	out := new(DeleteRoomResponse)
	err := c.cc.Invoke(ctx, RoomService_DeleteRoom_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RoomServiceServer is the server API for RoomService service.
// All implementations should embed UnimplementedRoomServiceServer
// for forward compatibility
type RoomServiceServer interface {
	CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error)
	GetRoomsByUser(context.Context, *GetRoomsByUserRequest) (*GetRoomsByUserResponse, error)
	GetUserIDsByRoom(context.Context, *GetUserIDsByRoomRequest) (*GetUserIDsByRoomResponse, error)
	InviteToRoom(context.Context, *InviteToRoomRequest) (*InviteToRoomResponse, error)
	EnterRoom(*EnterRoomRequest, RoomService_EnterRoomServer) error
	SendMessage(context.Context, *Message) (*SendMessageResponse, error)
	UpdateRoom(context.Context, *UpdateRoomRequest) (*UpdateRoomResponse, error)
	DeleteRoom(context.Context, *DeleteRoomRequest) (*DeleteRoomResponse, error)
}

// UnimplementedRoomServiceServer should be embedded to have forward compatible implementations.
type UnimplementedRoomServiceServer struct {
}

func (UnimplementedRoomServiceServer) CreateRoom(context.Context, *CreateRoomRequest) (*CreateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateRoom not implemented")
}
func (UnimplementedRoomServiceServer) GetRoomsByUser(context.Context, *GetRoomsByUserRequest) (*GetRoomsByUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRoomsByUser not implemented")
}
func (UnimplementedRoomServiceServer) GetUserIDsByRoom(context.Context, *GetUserIDsByRoomRequest) (*GetUserIDsByRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserIDsByRoom not implemented")
}
func (UnimplementedRoomServiceServer) InviteToRoom(context.Context, *InviteToRoomRequest) (*InviteToRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteToRoom not implemented")
}
func (UnimplementedRoomServiceServer) EnterRoom(*EnterRoomRequest, RoomService_EnterRoomServer) error {
	return status.Errorf(codes.Unimplemented, "method EnterRoom not implemented")
}
func (UnimplementedRoomServiceServer) SendMessage(context.Context, *Message) (*SendMessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedRoomServiceServer) UpdateRoom(context.Context, *UpdateRoomRequest) (*UpdateRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRoom not implemented")
}
func (UnimplementedRoomServiceServer) DeleteRoom(context.Context, *DeleteRoomRequest) (*DeleteRoomResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRoom not implemented")
}

// UnsafeRoomServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RoomServiceServer will
// result in compilation errors.
type UnsafeRoomServiceServer interface {
	mustEmbedUnimplementedRoomServiceServer()
}

func RegisterRoomServiceServer(s grpc.ServiceRegistrar, srv RoomServiceServer) {
	s.RegisterService(&RoomService_ServiceDesc, srv)
}

func _RoomService_CreateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).CreateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_CreateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).CreateRoom(ctx, req.(*CreateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_GetRoomsByUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRoomsByUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetRoomsByUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_GetRoomsByUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetRoomsByUser(ctx, req.(*GetRoomsByUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_GetUserIDsByRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserIDsByRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).GetUserIDsByRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_GetUserIDsByRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).GetUserIDsByRoom(ctx, req.(*GetUserIDsByRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_InviteToRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteToRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).InviteToRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_InviteToRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).InviteToRoom(ctx, req.(*InviteToRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_EnterRoom_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EnterRoomRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RoomServiceServer).EnterRoom(m, &roomServiceEnterRoomServer{stream})
}

type RoomService_EnterRoomServer interface {
	Send(*Message) error
	grpc.ServerStream
}

type roomServiceEnterRoomServer struct {
	grpc.ServerStream
}

func (x *roomServiceEnterRoomServer) Send(m *Message) error {
	return x.ServerStream.SendMsg(m)
}

func _RoomService_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Message)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).SendMessage(ctx, req.(*Message))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_UpdateRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).UpdateRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_UpdateRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).UpdateRoom(ctx, req.(*UpdateRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RoomService_DeleteRoom_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRoomRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RoomServiceServer).DeleteRoom(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RoomService_DeleteRoom_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RoomServiceServer).DeleteRoom(ctx, req.(*DeleteRoomRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RoomService_ServiceDesc is the grpc.ServiceDesc for RoomService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RoomService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "watch2gather.proto.roomapi.RoomService",
	HandlerType: (*RoomServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRoom",
			Handler:    _RoomService_CreateRoom_Handler,
		},
		{
			MethodName: "GetRoomsByUser",
			Handler:    _RoomService_GetRoomsByUser_Handler,
		},
		{
			MethodName: "GetUserIDsByRoom",
			Handler:    _RoomService_GetUserIDsByRoom_Handler,
		},
		{
			MethodName: "InviteToRoom",
			Handler:    _RoomService_InviteToRoom_Handler,
		},
		{
			MethodName: "SendMessage",
			Handler:    _RoomService_SendMessage_Handler,
		},
		{
			MethodName: "UpdateRoom",
			Handler:    _RoomService_UpdateRoom_Handler,
		},
		{
			MethodName: "DeleteRoom",
			Handler:    _RoomService_DeleteRoom_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "EnterRoom",
			Handler:       _RoomService_EnterRoom_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "room.proto",
}
