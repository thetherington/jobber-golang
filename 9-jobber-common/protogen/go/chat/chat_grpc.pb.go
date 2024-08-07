// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v5.26.1
// source: proto/chat/chat.proto

package chat

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	CreateMessage(ctx context.Context, in *MessageDocument, opts ...grpc.CallOption) (*MessageResponse, error)
	GetConversations(ctx context.Context, in *RequestMsgConversations, opts ...grpc.CallOption) (*ConversationsResponse, error)
	GetMessages(ctx context.Context, in *RequestMsgConversations, opts ...grpc.CallOption) (*MessagesResponse, error)
	GetConversationList(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*MessagesResponse, error)
	GetUserMessages(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*MessagesResponse, error)
	MarkMultipleMessages(ctx context.Context, in *RequestMsgConversations, opts ...grpc.CallOption) (*ResponseMessage, error)
	MarkSingleMessage(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*MessageResponse, error)
	UpdateOffer(ctx context.Context, in *UpdateOfferRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	Subscribe(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (ChatService_SubscribeClient, error)
	Unsubscribe(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) CreateMessage(ctx context.Context, in *MessageDocument, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/chat.ChatService/CreateMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetConversations(ctx context.Context, in *RequestMsgConversations, opts ...grpc.CallOption) (*ConversationsResponse, error) {
	out := new(ConversationsResponse)
	err := c.cc.Invoke(ctx, "/chat.ChatService/GetConversations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetMessages(ctx context.Context, in *RequestMsgConversations, opts ...grpc.CallOption) (*MessagesResponse, error) {
	out := new(MessagesResponse)
	err := c.cc.Invoke(ctx, "/chat.ChatService/GetMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetConversationList(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*MessagesResponse, error) {
	out := new(MessagesResponse)
	err := c.cc.Invoke(ctx, "/chat.ChatService/GetConversationList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetUserMessages(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*MessagesResponse, error) {
	out := new(MessagesResponse)
	err := c.cc.Invoke(ctx, "/chat.ChatService/GetUserMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) MarkMultipleMessages(ctx context.Context, in *RequestMsgConversations, opts ...grpc.CallOption) (*ResponseMessage, error) {
	out := new(ResponseMessage)
	err := c.cc.Invoke(ctx, "/chat.ChatService/MarkMultipleMessages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) MarkSingleMessage(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/chat.ChatService/MarkSingleMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) UpdateOffer(ctx context.Context, in *UpdateOfferRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/chat.ChatService/UpdateOffer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) Subscribe(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (ChatService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &ChatService_ServiceDesc.Streams[0], "/chat.ChatService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &chatServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ChatService_SubscribeClient interface {
	Recv() (*MessageResponse, error)
	grpc.ClientStream
}

type chatServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *chatServiceSubscribeClient) Recv() (*MessageResponse, error) {
	m := new(MessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *chatServiceClient) Unsubscribe(ctx context.Context, in *RequestWithParam, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/chat.ChatService/Unsubscribe", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	CreateMessage(context.Context, *MessageDocument) (*MessageResponse, error)
	GetConversations(context.Context, *RequestMsgConversations) (*ConversationsResponse, error)
	GetMessages(context.Context, *RequestMsgConversations) (*MessagesResponse, error)
	GetConversationList(context.Context, *RequestWithParam) (*MessagesResponse, error)
	GetUserMessages(context.Context, *RequestWithParam) (*MessagesResponse, error)
	MarkMultipleMessages(context.Context, *RequestMsgConversations) (*ResponseMessage, error)
	MarkSingleMessage(context.Context, *RequestWithParam) (*MessageResponse, error)
	UpdateOffer(context.Context, *UpdateOfferRequest) (*MessageResponse, error)
	Subscribe(*RequestWithParam, ChatService_SubscribeServer) error
	Unsubscribe(context.Context, *RequestWithParam) (*emptypb.Empty, error)
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) CreateMessage(context.Context, *MessageDocument) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMessage not implemented")
}
func (UnimplementedChatServiceServer) GetConversations(context.Context, *RequestMsgConversations) (*ConversationsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversations not implemented")
}
func (UnimplementedChatServiceServer) GetMessages(context.Context, *RequestMsgConversations) (*MessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}
func (UnimplementedChatServiceServer) GetConversationList(context.Context, *RequestWithParam) (*MessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversationList not implemented")
}
func (UnimplementedChatServiceServer) GetUserMessages(context.Context, *RequestWithParam) (*MessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserMessages not implemented")
}
func (UnimplementedChatServiceServer) MarkMultipleMessages(context.Context, *RequestMsgConversations) (*ResponseMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkMultipleMessages not implemented")
}
func (UnimplementedChatServiceServer) MarkSingleMessage(context.Context, *RequestWithParam) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkSingleMessage not implemented")
}
func (UnimplementedChatServiceServer) UpdateOffer(context.Context, *UpdateOfferRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateOffer not implemented")
}
func (UnimplementedChatServiceServer) Subscribe(*RequestWithParam, ChatService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedChatServiceServer) Unsubscribe(context.Context, *RequestWithParam) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Unsubscribe not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_CreateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageDocument)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).CreateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/CreateMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).CreateMessage(ctx, req.(*MessageDocument))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetConversations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMsgConversations)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetConversations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/GetConversations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetConversations(ctx, req.(*RequestMsgConversations))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMsgConversations)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/GetMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetMessages(ctx, req.(*RequestMsgConversations))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetConversationList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetConversationList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/GetConversationList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetConversationList(ctx, req.(*RequestWithParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetUserMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetUserMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/GetUserMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetUserMessages(ctx, req.(*RequestWithParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_MarkMultipleMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestMsgConversations)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).MarkMultipleMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/MarkMultipleMessages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).MarkMultipleMessages(ctx, req.(*RequestMsgConversations))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_MarkSingleMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).MarkSingleMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/MarkSingleMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).MarkSingleMessage(ctx, req.(*RequestWithParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_UpdateOffer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateOfferRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).UpdateOffer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/UpdateOffer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).UpdateOffer(ctx, req.(*UpdateOfferRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(RequestWithParam)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServiceServer).Subscribe(m, &chatServiceSubscribeServer{stream})
}

type ChatService_SubscribeServer interface {
	Send(*MessageResponse) error
	grpc.ServerStream
}

type chatServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *chatServiceSubscribeServer) Send(m *MessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ChatService_Unsubscribe_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestWithParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).Unsubscribe(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/chat.ChatService/Unsubscribe",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).Unsubscribe(ctx, req.(*RequestWithParam))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "chat.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMessage",
			Handler:    _ChatService_CreateMessage_Handler,
		},
		{
			MethodName: "GetConversations",
			Handler:    _ChatService_GetConversations_Handler,
		},
		{
			MethodName: "GetMessages",
			Handler:    _ChatService_GetMessages_Handler,
		},
		{
			MethodName: "GetConversationList",
			Handler:    _ChatService_GetConversationList_Handler,
		},
		{
			MethodName: "GetUserMessages",
			Handler:    _ChatService_GetUserMessages_Handler,
		},
		{
			MethodName: "MarkMultipleMessages",
			Handler:    _ChatService_MarkMultipleMessages_Handler,
		},
		{
			MethodName: "MarkSingleMessage",
			Handler:    _ChatService_MarkSingleMessage_Handler,
		},
		{
			MethodName: "UpdateOffer",
			Handler:    _ChatService_UpdateOffer_Handler,
		},
		{
			MethodName: "Unsubscribe",
			Handler:    _ChatService_Unsubscribe_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _ChatService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/chat/chat.proto",
}
