// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: movie.proto

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
	MovieService_GetAllMovies_FullMethodName   = "/watch2gather.proto.movieapi.MovieService/GetAllMovies"
	MovieService_GetMovie_FullMethodName       = "/watch2gather.proto.movieapi.MovieService/GetMovie"
	MovieService_GetMoviePoster_FullMethodName = "/watch2gather.proto.movieapi.MovieService/GetMoviePoster"
)

// MovieServiceClient is the client API for MovieService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MovieServiceClient interface {
	GetAllMovies(ctx context.Context, in *GetAllMoviesRequest, opts ...grpc.CallOption) (*GetAllMoviesResponse, error)
	GetMovie(ctx context.Context, in *GetMovieRequest, opts ...grpc.CallOption) (*Movie, error)
	GetMoviePoster(ctx context.Context, in *GetMoviePosterRequest, opts ...grpc.CallOption) (*GetMoviePosterResponse, error)
}

type movieServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMovieServiceClient(cc grpc.ClientConnInterface) MovieServiceClient {
	return &movieServiceClient{cc}
}

func (c *movieServiceClient) GetAllMovies(ctx context.Context, in *GetAllMoviesRequest, opts ...grpc.CallOption) (*GetAllMoviesResponse, error) {
	out := new(GetAllMoviesResponse)
	err := c.cc.Invoke(ctx, MovieService_GetAllMovies_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *movieServiceClient) GetMovie(ctx context.Context, in *GetMovieRequest, opts ...grpc.CallOption) (*Movie, error) {
	out := new(Movie)
	err := c.cc.Invoke(ctx, MovieService_GetMovie_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *movieServiceClient) GetMoviePoster(ctx context.Context, in *GetMoviePosterRequest, opts ...grpc.CallOption) (*GetMoviePosterResponse, error) {
	out := new(GetMoviePosterResponse)
	err := c.cc.Invoke(ctx, MovieService_GetMoviePoster_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MovieServiceServer is the server API for MovieService service.
// All implementations should embed UnimplementedMovieServiceServer
// for forward compatibility
type MovieServiceServer interface {
	GetAllMovies(context.Context, *GetAllMoviesRequest) (*GetAllMoviesResponse, error)
	GetMovie(context.Context, *GetMovieRequest) (*Movie, error)
	GetMoviePoster(context.Context, *GetMoviePosterRequest) (*GetMoviePosterResponse, error)
}

// UnimplementedMovieServiceServer should be embedded to have forward compatible implementations.
type UnimplementedMovieServiceServer struct {
}

func (UnimplementedMovieServiceServer) GetAllMovies(context.Context, *GetAllMoviesRequest) (*GetAllMoviesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllMovies not implemented")
}
func (UnimplementedMovieServiceServer) GetMovie(context.Context, *GetMovieRequest) (*Movie, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMovie not implemented")
}
func (UnimplementedMovieServiceServer) GetMoviePoster(context.Context, *GetMoviePosterRequest) (*GetMoviePosterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMoviePoster not implemented")
}

// UnsafeMovieServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MovieServiceServer will
// result in compilation errors.
type UnsafeMovieServiceServer interface {
	mustEmbedUnimplementedMovieServiceServer()
}

func RegisterMovieServiceServer(s grpc.ServiceRegistrar, srv MovieServiceServer) {
	s.RegisterService(&MovieService_ServiceDesc, srv)
}

func _MovieService_GetAllMovies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllMoviesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovieServiceServer).GetAllMovies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MovieService_GetAllMovies_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovieServiceServer).GetAllMovies(ctx, req.(*GetAllMoviesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MovieService_GetMovie_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMovieRequest)
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
		return srv.(MovieServiceServer).GetMovie(ctx, req.(*GetMovieRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MovieService_GetMoviePoster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMoviePosterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovieServiceServer).GetMoviePoster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MovieService_GetMoviePoster_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovieServiceServer).GetMoviePoster(ctx, req.(*GetMoviePosterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MovieService_ServiceDesc is the grpc.ServiceDesc for MovieService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MovieService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "watch2gather.proto.movieapi.MovieService",
	HandlerType: (*MovieServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAllMovies",
			Handler:    _MovieService_GetAllMovies_Handler,
		},
		{
			MethodName: "GetMovie",
			Handler:    _MovieService_GetMovie_Handler,
		},
		{
			MethodName: "GetMoviePoster",
			Handler:    _MovieService_GetMoviePoster_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "movie.proto",
}
