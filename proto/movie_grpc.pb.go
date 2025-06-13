// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: movie.proto

package movie_grpc

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
	MovieService_CreateMovie_FullMethodName = "/proto.MovieService/CreateMovie"
	MovieService_GetMovie_FullMethodName    = "/proto.MovieService/GetMovie"
	MovieService_GetMovies_FullMethodName   = "/proto.MovieService/GetMovies"
	MovieService_UpdateMovie_FullMethodName = "/proto.MovieService/UpdateMovie"
	MovieService_DeleteMovie_FullMethodName = "/proto.MovieService/DeleteMovie"
)

// MovieServiceClient is the client API for MovieService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MovieServiceClient interface {
	CreateMovie(ctx context.Context, in *CreateMovieRequest, opts ...grpc.CallOption) (*CreateMovieResponse, error)
	GetMovie(ctx context.Context, in *ReadMovieRequest, opts ...grpc.CallOption) (*ReadMovieResponse, error)
	GetMovies(ctx context.Context, in *ReadMoviesRequest, opts ...grpc.CallOption) (*ReadMoviesResponse, error)
	UpdateMovie(ctx context.Context, in *UpdateMovieRequest, opts ...grpc.CallOption) (*UpdateMovieResponse, error)
	DeleteMovie(ctx context.Context, in *DeleteMovieRequest, opts ...grpc.CallOption) (*DeleteMovieResponse, error)
}

type movieServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMovieServiceClient(cc grpc.ClientConnInterface) MovieServiceClient {
	return &movieServiceClient{cc}
}

func (c *movieServiceClient) CreateMovie(ctx context.Context, in *CreateMovieRequest, opts ...grpc.CallOption) (*CreateMovieResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateMovieResponse)
	err := c.cc.Invoke(ctx, MovieService_CreateMovie_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *movieServiceClient) GetMovie(ctx context.Context, in *ReadMovieRequest, opts ...grpc.CallOption) (*ReadMovieResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadMovieResponse)
	err := c.cc.Invoke(ctx, MovieService_GetMovie_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *movieServiceClient) GetMovies(ctx context.Context, in *ReadMoviesRequest, opts ...grpc.CallOption) (*ReadMoviesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadMoviesResponse)
	err := c.cc.Invoke(ctx, MovieService_GetMovies_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *movieServiceClient) UpdateMovie(ctx context.Context, in *UpdateMovieRequest, opts ...grpc.CallOption) (*UpdateMovieResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMovieResponse)
	err := c.cc.Invoke(ctx, MovieService_UpdateMovie_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *movieServiceClient) DeleteMovie(ctx context.Context, in *DeleteMovieRequest, opts ...grpc.CallOption) (*DeleteMovieResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteMovieResponse)
	err := c.cc.Invoke(ctx, MovieService_DeleteMovie_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MovieServiceServer is the server API for MovieService service.
// All implementations must embed UnimplementedMovieServiceServer
// for forward compatibility.
type MovieServiceServer interface {
	CreateMovie(context.Context, *CreateMovieRequest) (*CreateMovieResponse, error)
	GetMovie(context.Context, *ReadMovieRequest) (*ReadMovieResponse, error)
	GetMovies(context.Context, *ReadMoviesRequest) (*ReadMoviesResponse, error)
	UpdateMovie(context.Context, *UpdateMovieRequest) (*UpdateMovieResponse, error)
	DeleteMovie(context.Context, *DeleteMovieRequest) (*DeleteMovieResponse, error)
	mustEmbedUnimplementedMovieServiceServer()
}

// UnimplementedMovieServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMovieServiceServer struct{}

func (UnimplementedMovieServiceServer) CreateMovie(context.Context, *CreateMovieRequest) (*CreateMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMovie not implemented")
}
func (UnimplementedMovieServiceServer) GetMovie(context.Context, *ReadMovieRequest) (*ReadMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovie not implemented")
}
func (UnimplementedMovieServiceServer) GetMovies(context.Context, *ReadMoviesRequest) (*ReadMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovies not implemented")
}
func (UnimplementedMovieServiceServer) UpdateMovie(context.Context, *UpdateMovieRequest) (*UpdateMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMovie not implemented")
}
func (UnimplementedMovieServiceServer) DeleteMovie(context.Context, *DeleteMovieRequest) (*DeleteMovieResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMovie not implemented")
}
func (UnimplementedMovieServiceServer) mustEmbedUnimplementedMovieServiceServer() {}
func (UnimplementedMovieServiceServer) testEmbeddedByValue()                      {}

// UnsafeMovieServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MovieServiceServer will
// result in compilation errors.
type UnsafeMovieServiceServer interface {
	mustEmbedUnimplementedMovieServiceServer()
}

func RegisterMovieServiceServer(s grpc.ServiceRegistrar, srv MovieServiceServer) {
	// If the following call pancis, it indicates UnimplementedMovieServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&MovieService_ServiceDesc, srv)
}

func _MovieService_CreateMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMovieRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovieServiceServer).CreateMovie(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MovieService_CreateMovie_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovieServiceServer).CreateMovie(ctx, req.(*CreateMovieRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MovieService_GetMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadMovieRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovieServiceServer).GetMovie(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MovieService_GetMovie_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovieServiceServer).GetMovie(ctx, req.(*ReadMovieRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MovieService_GetMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadMoviesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovieServiceServer).GetMovies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MovieService_GetMovies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovieServiceServer).GetMovies(ctx, req.(*ReadMoviesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MovieService_UpdateMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMovieRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovieServiceServer).UpdateMovie(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MovieService_UpdateMovie_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovieServiceServer).UpdateMovie(ctx, req.(*UpdateMovieRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MovieService_DeleteMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMovieRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovieServiceServer).DeleteMovie(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MovieService_DeleteMovie_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovieServiceServer).DeleteMovie(ctx, req.(*DeleteMovieRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MovieService_ServiceDesc is the grpc.ServiceDesc for MovieService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MovieService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.MovieService",
	HandlerType: (*MovieServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMovie",
			Handler:    _MovieService_CreateMovie_Handler,
		},
		{
			MethodName: "GetMovie",
			Handler:    _MovieService_GetMovie_Handler,
		},
		{
			MethodName: "GetMovies",
			Handler:    _MovieService_GetMovies_Handler,
		},
		{
			MethodName: "UpdateMovie",
			Handler:    _MovieService_UpdateMovie_Handler,
		},
		{
			MethodName: "DeleteMovie",
			Handler:    _MovieService_DeleteMovie_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "movie.proto",
}
