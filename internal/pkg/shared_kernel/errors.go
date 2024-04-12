package sharedkernel

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrServer          = status.Error(codes.Internal, "Server error")
	ErrMissingMetadata = status.Error(codes.InvalidArgument, "missing metadata")
	ErrInvalidToken    = status.Error(codes.Unauthenticated, "invalid token")
	ErrUsersNotFound   = status.Error(codes.InvalidArgument, "users not found")
)
