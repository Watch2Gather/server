package router

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Watch2Gather/server/cmd/user/config"
	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/user/usecases/users"
	"github.com/Watch2Gather/server/proto/gen"
)

type userInfoGRPCServer struct {
	cfg *config.Config
	uc  users.UseCase
}


var _ gen.UserInfoServiceServer = (*userInfoGRPCServer)(nil)

var UserInfoGRPCServerSet = wire.NewSet(NewGRPCUserInfoServer)

func NewGRPCUserInfoServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc users.UseCase,
) gen.UserInfoServiceServer {
	svc := userInfoGRPCServer{
		cfg: cfg,
		uc:  uc,
	}

	gen.RegisterUserInfoServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (g *userInfoGRPCServer) GetUserInfo(ctx context.Context, req *gen.GetUserInfoRequest) (*gen.GetUserInfoResponse, error) {
	slog.Info("gRPC client", "http_method", "GET", "http_name", "GetItemTypes")

	id, err := uuid.Parse(req.GetId())
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	user, err := g.uc.GetUserData(ctx, id)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.GetUserData"))
		return nil, sharedkernel.ErrServer
	}

	return &gen.GetUserInfoResponse{
		User: &gen.UserInfo{
			Username: user.Username,
			Avatar: user.Avatar,
			Id: req.GetId(),
		},
	}, nil
}

