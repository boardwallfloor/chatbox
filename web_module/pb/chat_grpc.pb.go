// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: chat.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	// Check for unsent messages
	CheckNewMessages(ctx context.Context, in *timestamppb.Timestamp, opts ...grpc.CallOption) (MessageService_CheckNewMessagesClient, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) CheckNewMessages(ctx context.Context, in *timestamppb.Timestamp, opts ...grpc.CallOption) (MessageService_CheckNewMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &MessageService_ServiceDesc.Streams[0], "/myapp.MessageService/CheckNewMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &messageServiceCheckNewMessagesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MessageService_CheckNewMessagesClient interface {
	Recv() (*ChatMessage, error)
	grpc.ClientStream
}

type messageServiceCheckNewMessagesClient struct {
	grpc.ClientStream
}

func (x *messageServiceCheckNewMessagesClient) Recv() (*ChatMessage, error) {
	m := new(ChatMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	// Check for unsent messages
	CheckNewMessages(*timestamppb.Timestamp, MessageService_CheckNewMessagesServer) error
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) CheckNewMessages(*timestamppb.Timestamp, MessageService_CheckNewMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method CheckNewMessages not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_CheckNewMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(timestamppb.Timestamp)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MessageServiceServer).CheckNewMessages(m, &messageServiceCheckNewMessagesServer{stream})
}

type MessageService_CheckNewMessagesServer interface {
	Send(*ChatMessage) error
	grpc.ServerStream
}

type messageServiceCheckNewMessagesServer struct {
	grpc.ServerStream
}

func (x *messageServiceCheckNewMessagesServer) Send(m *ChatMessage) error {
	return x.ServerStream.SendMsg(m)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "myapp.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "CheckNewMessages",
			Handler:       _MessageService_CheckNewMessages_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "chat.proto",
}
