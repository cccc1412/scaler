// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package modelproto

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

// ModelClient is the client API for Model service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ModelClient interface {
	Predict(ctx context.Context, in *PredictRequest, opts ...grpc.CallOption) (*PredictResponse, error)
}

type modelClient struct {
	cc grpc.ClientConnInterface
}

func NewModelClient(cc grpc.ClientConnInterface) ModelClient {
	return &modelClient{cc}
}

func (c *modelClient) Predict(ctx context.Context, in *PredictRequest, opts ...grpc.CallOption) (*PredictResponse, error) {
	out := new(PredictResponse)
	err := c.cc.Invoke(ctx, "/serverless.model.Model/Predict", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ModelServer is the server API for Model service.
// All implementations must embed UnimplementedModelServer
// for forward compatibility
type ModelServer interface {
	Predict(context.Context, *PredictRequest) (*PredictResponse, error)
	mustEmbedUnimplementedModelServer()
}

// UnimplementedModelServer must be embedded to have forward compatible implementations.
type UnimplementedModelServer struct {
}

func (UnimplementedModelServer) Predict(context.Context, *PredictRequest) (*PredictResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Predict not implemented")
}
func (UnimplementedModelServer) mustEmbedUnimplementedModelServer() {}

// UnsafeModelServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ModelServer will
// result in compilation errors.
type UnsafeModelServer interface {
	mustEmbedUnimplementedModelServer()
}

func RegisterModelServer(s grpc.ServiceRegistrar, srv ModelServer) {
	s.RegisterService(&Model_ServiceDesc, srv)
}

func _Model_Predict_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PredictRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModelServer).Predict(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/serverless.model.Model/Predict",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModelServer).Predict(ctx, req.(*PredictRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Model_ServiceDesc is the grpc.ServiceDesc for Model service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Model_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "serverless.model.Model",
	HandlerType: (*ModelServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Predict",
			Handler:    _Model_Predict_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "model_proto/model.proto",
}
