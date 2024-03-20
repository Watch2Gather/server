package router

import (
	"context"
	"log/slog"

	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Watch2Gather/server/cmd/user/config"
	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
	"github.com/Watch2Gather/server/internal/user/domain"
	"github.com/Watch2Gather/server/internal/user/usecases/users"
	"github.com/Watch2Gather/server/proto/gen"
)

type userGRPCServer struct {
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

func (g *userGRPCServer) RegisterUser(ctx context.Context, req *gen.RegisterUserRequest) (*gen.RegisterUserResponse, error) {
	slog.Info("POST: RegisterUser")

	if len(req.Username) < 1 || len(req.Username) > 36{
		return &gen.RegisterUserResponse{}, errors.New("username must be between 1 and 36 characters")
	}

	if len(req.Email) < 3 || len(req.Email) > 254{
		return nil, errors.New("email must be between 3 and 254 characters")
	}

	err := g.uc.Register(ctx, &domain.RegisterModel{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		slog.Error("Caught error", "error", errors.Wrap(err, "uc.Register"))
		return nil, sharedkernel.ErrServer
	}

	res := gen.RegisterUserResponse{}

	return &res, nil
}

func (g *userGRPCServer) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	slog.Info("POST: RegisterUser")

	res := gen.LoginUserResponse{}
	token, err := g.uc.Login(ctx, &domain.LoginModel{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		slog.Debug("error type", "err: ", err.Error())
		if errors.Is(err, domain.ErrUnauthorized){
			return nil, errors.New("Invalid username or password")
		}
		return nil, errors.Wrap(err, "uc.Login")
	}

	res.AccessToken = token.AccessToken
	res.RefreshToken = token.RefreshToken

	return &res, nil
}

func (userGRPCServer) ChangeUserData(_ context.Context, _ *gen.ChangeUserDataRequest) (_ *gen.ChangeUserDataResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (userGRPCServer) ChangePassword(_ context.Context, _ *gen.ChangePasswordRequest) (_ *gen.ChangePasswordResponse, _ error) {
	panic("not implemented") // TODO: Implement
}

func (userGRPCServer) RefreshToken(_ context.Context, _ *gen.RefreshTokenRequest) (_ *gen.RefreshTokenResponse, _ error) {
	panic("not implemented") // TODO: Implement
}
