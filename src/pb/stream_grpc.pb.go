// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: stream.proto

package pb

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
	StreamService_FocusPointStream_FullMethodName   = "/pb.StreamService/FocusPointStream"
	StreamService_ImageToVideoStream_FullMethodName = "/pb.StreamService/ImageToVideoStream"
)

// StreamServiceClient is the client API for StreamService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StreamServiceClient interface {
	// Focus point method
	FocusPointStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_FocusPointStreamClient, error)
	// image to video method
	ImageToVideoStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_ImageToVideoStreamClient, error)
}

type streamServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStreamServiceClient(cc grpc.ClientConnInterface) StreamServiceClient {
	return &streamServiceClient{cc}
}

func (c *streamServiceClient) FocusPointStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_FocusPointStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[0], StreamService_FocusPointStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceFocusPointStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamService_FocusPointStreamClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type streamServiceFocusPointStreamClient struct {
	grpc.ClientStream
}

func (x *streamServiceFocusPointStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *streamServiceClient) ImageToVideoStream(ctx context.Context, in *Request, opts ...grpc.CallOption) (StreamService_ImageToVideoStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &StreamService_ServiceDesc.Streams[1], StreamService_ImageToVideoStream_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &streamServiceImageToVideoStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type StreamService_ImageToVideoStreamClient interface {
	Recv() (*Response, error)
	grpc.ClientStream
}

type streamServiceImageToVideoStreamClient struct {
	grpc.ClientStream
}

func (x *streamServiceImageToVideoStreamClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// StreamServiceServer is the server API for StreamService service.
// All implementations must embed UnimplementedStreamServiceServer
// for forward compatibility
type StreamServiceServer interface {
	// Focus point method
	FocusPointStream(*Request, StreamService_FocusPointStreamServer) error
	// image to video method
	ImageToVideoStream(*Request, StreamService_ImageToVideoStreamServer) error
	mustEmbedUnimplementedStreamServiceServer()
}

// UnimplementedStreamServiceServer must be embedded to have forward compatible implementations.
type UnimplementedStreamServiceServer struct {
}

func (UnimplementedStreamServiceServer) FocusPointStream(*Request, StreamService_FocusPointStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method FocusPointStream not implemented")
}
func (UnimplementedStreamServiceServer) ImageToVideoStream(*Request, StreamService_ImageToVideoStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method ImageToVideoStream not implemented")
}
func (UnimplementedStreamServiceServer) mustEmbedUnimplementedStreamServiceServer() {}

// UnsafeStreamServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StreamServiceServer will
// result in compilation errors.
type UnsafeStreamServiceServer interface {
	mustEmbedUnimplementedStreamServiceServer()
}

func RegisterStreamServiceServer(s grpc.ServiceRegistrar, srv StreamServiceServer) {
	s.RegisterService(&StreamService_ServiceDesc, srv)
}

func _StreamService_FocusPointStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).FocusPointStream(m, &streamServiceFocusPointStreamServer{stream})
}

type StreamService_FocusPointStreamServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type streamServiceFocusPointStreamServer struct {
	grpc.ServerStream
}

func (x *streamServiceFocusPointStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func _StreamService_ImageToVideoStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Request)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(StreamServiceServer).ImageToVideoStream(m, &streamServiceImageToVideoStreamServer{stream})
}

type StreamService_ImageToVideoStreamServer interface {
	Send(*Response) error
	grpc.ServerStream
}

type streamServiceImageToVideoStreamServer struct {
	grpc.ServerStream
}

func (x *streamServiceImageToVideoStreamServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

// StreamService_ServiceDesc is the grpc.ServiceDesc for StreamService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StreamService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.StreamService",
	HandlerType: (*StreamServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "FocusPointStream",
			Handler:       _StreamService_FocusPointStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "ImageToVideoStream",
			Handler:       _StreamService_ImageToVideoStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "stream.proto",
}