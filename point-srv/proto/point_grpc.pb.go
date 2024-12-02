// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package proto

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

// PointsClient is the client API for Points service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PointsClient interface {
	// 增加并冻结积分
	AddAndFreezePoints(ctx context.Context, in *AddAndFreezePointsRequest, opts ...grpc.CallOption) (*AddAndFreezePointsResponse, error)
	// 解冻积分
	UnfreezePoints(ctx context.Context, in *UnfreezePointsRequest, opts ...grpc.CallOption) (*UnfreezePointsResponse, error)
	// 获取积分明细
	GetPointsDetails(ctx context.Context, in *GetPointsDetailsRequest, opts ...grpc.CallOption) (*GetPointsDetailsResponse, error)
}

type pointsClient struct {
	cc grpc.ClientConnInterface
}

func NewPointsClient(cc grpc.ClientConnInterface) PointsClient {
	return &pointsClient{cc}
}

func (c *pointsClient) AddAndFreezePoints(ctx context.Context, in *AddAndFreezePointsRequest, opts ...grpc.CallOption) (*AddAndFreezePointsResponse, error) {
	out := new(AddAndFreezePointsResponse)
	err := c.cc.Invoke(ctx, "/Points/AddAndFreezePoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) UnfreezePoints(ctx context.Context, in *UnfreezePointsRequest, opts ...grpc.CallOption) (*UnfreezePointsResponse, error) {
	out := new(UnfreezePointsResponse)
	err := c.cc.Invoke(ctx, "/Points/UnfreezePoints", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pointsClient) GetPointsDetails(ctx context.Context, in *GetPointsDetailsRequest, opts ...grpc.CallOption) (*GetPointsDetailsResponse, error) {
	out := new(GetPointsDetailsResponse)
	err := c.cc.Invoke(ctx, "/Points/GetPointsDetails", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PointsServer is the server API for Points service.
// All implementations must embed UnimplementedPointsServer
// for forward compatibility
type PointsServer interface {
	// 增加并冻结积分
	AddAndFreezePoints(context.Context, *AddAndFreezePointsRequest) (*AddAndFreezePointsResponse, error)
	// 解冻积分
	UnfreezePoints(context.Context, *UnfreezePointsRequest) (*UnfreezePointsResponse, error)
	// 获取积分明细
	GetPointsDetails(context.Context, *GetPointsDetailsRequest) (*GetPointsDetailsResponse, error)
	mustEmbedUnimplementedPointsServer()
}

// UnimplementedPointsServer must be embedded to have forward compatible implementations.
type UnimplementedPointsServer struct {
}

func (UnimplementedPointsServer) AddAndFreezePoints(context.Context, *AddAndFreezePointsRequest) (*AddAndFreezePointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAndFreezePoints not implemented")
}
func (UnimplementedPointsServer) UnfreezePoints(context.Context, *UnfreezePointsRequest) (*UnfreezePointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnfreezePoints not implemented")
}
func (UnimplementedPointsServer) GetPointsDetails(context.Context, *GetPointsDetailsRequest) (*GetPointsDetailsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPointsDetails not implemented")
}
func (UnimplementedPointsServer) mustEmbedUnimplementedPointsServer() {}

// UnsafePointsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PointsServer will
// result in compilation errors.
type UnsafePointsServer interface {
	mustEmbedUnimplementedPointsServer()
}

func RegisterPointsServer(s grpc.ServiceRegistrar, srv PointsServer) {
	s.RegisterService(&Points_ServiceDesc, srv)
}

func _Points_AddAndFreezePoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAndFreezePointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PointsServer).AddAndFreezePoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Points/AddAndFreezePoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PointsServer).AddAndFreezePoints(ctx, req.(*AddAndFreezePointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Points_UnfreezePoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnfreezePointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PointsServer).UnfreezePoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Points/UnfreezePoints",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PointsServer).UnfreezePoints(ctx, req.(*UnfreezePointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Points_GetPointsDetails_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPointsDetailsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PointsServer).GetPointsDetails(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Points/GetPointsDetails",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PointsServer).GetPointsDetails(ctx, req.(*GetPointsDetailsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Points_ServiceDesc is the grpc.ServiceDesc for Points service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Points_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Points",
	HandlerType: (*PointsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddAndFreezePoints",
			Handler:    _Points_AddAndFreezePoints_Handler,
		},
		{
			MethodName: "UnfreezePoints",
			Handler:    _Points_UnfreezePoints_Handler,
		},
		{
			MethodName: "GetPointsDetails",
			Handler:    _Points_GetPointsDetails_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "point.proto",
}