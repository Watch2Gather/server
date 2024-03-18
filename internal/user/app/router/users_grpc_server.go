package router

import (
	"context"

	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Watch2Gather/server/cmd/user/config"
	"github.com/Watch2Gather/server/internal/user/usecases/users"
	"github.com/Watch2Gather/server/proto/gen"
)

type userGRPCServer struct {
	gen.UnimplementedUserServiceServer
	cfg *config.Config
	uc  users.UseCase
}

var _ gen.UserServiceServer = (*userGRPCServer)(nil)

var UserGRPCServerSet = wire.NewSet(NewGRPCUsersServer)

func NewGRPCUsersServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc users.UseCase,
) gen.UserServiceServer {
	svc := userGRPCServer{
		cfg: cfg,
		uc:  uc,
	}

	gen.RegisterUserServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

// TODO: Update Proto file

func (userGRPCServer) RegisterUser(_ context.Context, _ *gen.RegisterUserRequest) (_ *gen.RegisterUserResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (userGRPCServer) LoginUser(_ context.Context, _ *gen.LoginUserRequest) (_ *gen.LoginUserResponse, _ error) {
	panic("not implemented") // TODO: Implement
}
