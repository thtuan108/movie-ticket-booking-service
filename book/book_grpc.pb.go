// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package book

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

// BookerClient is the client API for Booker service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BookerClient interface {
	ListMovie(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*Options, error)
	ListGenres(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*Options, error)
	ListShowTimes(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*Options, error)
	MoviesByTitle(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (*Movies, error)
	MoviesByGenre(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (*Movies, error)
	MoviesByTime(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (*Movies, error)
	BuyTicket(ctx context.Context, in *BuyRequest, opts ...grpc.CallOption) (*NoParam, error)
}

type bookerClient struct {
	cc grpc.ClientConnInterface
}

func NewBookerClient(cc grpc.ClientConnInterface) BookerClient {
	return &bookerClient{cc}
}

func (c *bookerClient) ListMovie(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*Options, error) {
	out := new(Options)
	err := c.cc.Invoke(ctx, "/book.Booker/ListMovie", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookerClient) ListGenres(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*Options, error) {
	out := new(Options)
	err := c.cc.Invoke(ctx, "/book.Booker/ListGenres", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookerClient) ListShowTimes(ctx context.Context, in *NoParam, opts ...grpc.CallOption) (*Options, error) {
	out := new(Options)
	err := c.cc.Invoke(ctx, "/book.Booker/ListShowTimes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookerClient) MoviesByTitle(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (*Movies, error) {
	out := new(Movies)
	err := c.cc.Invoke(ctx, "/book.Booker/MoviesByTitle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookerClient) MoviesByGenre(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (*Movies, error) {
	out := new(Movies)
	err := c.cc.Invoke(ctx, "/book.Booker/MoviesByGenre", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookerClient) MoviesByTime(ctx context.Context, in *QueryParam, opts ...grpc.CallOption) (*Movies, error) {
	out := new(Movies)
	err := c.cc.Invoke(ctx, "/book.Booker/MoviesByTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bookerClient) BuyTicket(ctx context.Context, in *BuyRequest, opts ...grpc.CallOption) (*NoParam, error) {
	out := new(NoParam)
	err := c.cc.Invoke(ctx, "/book.Booker/BuyTicket", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BookerServer is the server API for Booker service.
// All implementations must embed UnimplementedBookerServer
// for forward compatibility
type BookerServer interface {
	ListMovie(context.Context, *NoParam) (*Options, error)
	ListGenres(context.Context, *NoParam) (*Options, error)
	ListShowTimes(context.Context, *NoParam) (*Options, error)
	MoviesByTitle(context.Context, *QueryParam) (*Movies, error)
	MoviesByGenre(context.Context, *QueryParam) (*Movies, error)
	MoviesByTime(context.Context, *QueryParam) (*Movies, error)
	BuyTicket(context.Context, *BuyRequest) (*NoParam, error)
	mustEmbedUnimplementedBookerServer()
}

// UnimplementedBookerServer must be embedded to have forward compatible implementations.
type UnimplementedBookerServer struct {
}

func (UnimplementedBookerServer) ListMovie(context.Context, *NoParam) (*Options, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMovie not implemented")
}
func (UnimplementedBookerServer) ListGenres(context.Context, *NoParam) (*Options, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListGenres not implemented")
}
func (UnimplementedBookerServer) ListShowTimes(context.Context, *NoParam) (*Options, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListShowTimes not implemented")
}
func (UnimplementedBookerServer) MoviesByTitle(context.Context, *QueryParam) (*Movies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoviesByTitle not implemented")
}
func (UnimplementedBookerServer) MoviesByGenre(context.Context, *QueryParam) (*Movies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoviesByGenre not implemented")
}
func (UnimplementedBookerServer) MoviesByTime(context.Context, *QueryParam) (*Movies, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoviesByTime not implemented")
}
func (UnimplementedBookerServer) BuyTicket(context.Context, *BuyRequest) (*NoParam, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BuyTicket not implemented")
}
func (UnimplementedBookerServer) mustEmbedUnimplementedBookerServer() {}

// UnsafeBookerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BookerServer will
// result in compilation errors.
type UnsafeBookerServer interface {
	mustEmbedUnimplementedBookerServer()
}

func RegisterBookerServer(s grpc.ServiceRegistrar, srv BookerServer) {
	s.RegisterService(&Booker_ServiceDesc, srv)
}

func _Booker_ListMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookerServer).ListMovie(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.Booker/ListMovie",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookerServer).ListMovie(ctx, req.(*NoParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Booker_ListGenres_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookerServer).ListGenres(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.Booker/ListGenres",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookerServer).ListGenres(ctx, req.(*NoParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Booker_ListShowTimes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NoParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookerServer).ListShowTimes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.Booker/ListShowTimes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookerServer).ListShowTimes(ctx, req.(*NoParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Booker_MoviesByTitle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookerServer).MoviesByTitle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.Booker/MoviesByTitle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookerServer).MoviesByTitle(ctx, req.(*QueryParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Booker_MoviesByGenre_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookerServer).MoviesByGenre(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.Booker/MoviesByGenre",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookerServer).MoviesByGenre(ctx, req.(*QueryParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Booker_MoviesByTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParam)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookerServer).MoviesByTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.Booker/MoviesByTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookerServer).MoviesByTime(ctx, req.(*QueryParam))
	}
	return interceptor(ctx, in, info, handler)
}

func _Booker_BuyTicket_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BuyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BookerServer).BuyTicket(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/book.Booker/BuyTicket",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BookerServer).BuyTicket(ctx, req.(*BuyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Booker_ServiceDesc is the grpc.ServiceDesc for Booker service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Booker_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "book.Booker",
	HandlerType: (*BookerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListMovie",
			Handler:    _Booker_ListMovie_Handler,
		},
		{
			MethodName: "ListGenres",
			Handler:    _Booker_ListGenres_Handler,
		},
		{
			MethodName: "ListShowTimes",
			Handler:    _Booker_ListShowTimes_Handler,
		},
		{
			MethodName: "MoviesByTitle",
			Handler:    _Booker_MoviesByTitle_Handler,
		},
		{
			MethodName: "MoviesByGenre",
			Handler:    _Booker_MoviesByGenre_Handler,
		},
		{
			MethodName: "MoviesByTime",
			Handler:    _Booker_MoviesByTime_Handler,
		},
		{
			MethodName: "BuyTicket",
			Handler:    _Booker_BuyTicket_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "book/book.proto",
}
