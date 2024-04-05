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

	if len(req.GetUsername()) < 1 || len(req.GetUsername()) > 36 {
		return nil, errors.New("username must be between 1 and 36 characters")
	}

	if len(req.GetEmail()) < 3 || len(req.GetEmail()) > 254 {
		return nil, errors.New("email must be between 3 and 254 characters")
	}

	if len(req.GetPassword()) < 8 {
		return nil, errors.New("Password must be at least 8 characters")
	}

	err := g.uc.Register(ctx, &domain.RegisterModel{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	})
	if err != nil {
		slog.Error("Caught error", "error", errors.Wrap(err, "uc.Register"))
		return nil, sharedkernel.ErrServer
	}

	res := gen.RegisterUserResponse{}

	return &res, nil
}

func (g *userGRPCServer) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	slog.Info("POST: LoginUser")

	if len(req.GetUsername()) < 1 || len(req.GetUsername()) > 36 {
		return nil, errors.New("username must be between 1 and 36 characters")
	}

	if len(req.GetPassword()) < 8 {
		return nil, errors.New("Password must be at least 8 characters")
	}
	res := gen.LoginUserResponse{}
	token, err := g.uc.Login(ctx, &domain.LoginModel{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			return nil, errors.New("Invalid username or password")
		}
		return nil, errors.Wrap(err, "uc.Login")
	}

	res.AccessToken = token.AccessToken
	res.RefreshToken = token.RefreshToken

	return &res, nil
}

func (g *userGRPCServer) ChangeUserData(ctx context.Context, req *gen.ChangeUserDataRequest) (*gen.ChangeUserDataResponse, error) {
	slog.Info("PUT: ChangeUserData")

	if len(req.GetUsername()) < 1 || len(req.GetUsername()) > 36 {
		return nil, errors.New("username must be between 1 and 36 characters")
	}

	if len(req.GetEmail()) < 3 || len(req.GetEmail()) > 254 {
		return nil, errors.New("email must be between 3 and 254 characters")
	}


	tokenString, err := sharedkernel.GetToken(ctx)
	if err != nil{
		return nil, errors.Wrap(err, "sharedkernel.GetToken")
	}

	tokenClaims, err := sharedkernel.ParseToken(ctx, tokenString)
	if err != nil{
		return nil, errors.Wrap(err, "sharedkernel.ParseToken")
	}

	id, err := uuid.Parse(tokenClaims.Id)
  if err != nil{
    return nil, errors.Wrap(err, "uuid.Parse")
  }

	user, err := g.uc.ChangeUserData(ctx, &domain.User{
		ID:       id,
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Avatar:   req.GetAvatar(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "uc.ChangeUserData")
	}

	res := gen.ChangeUserDataResponse{
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.Avatar,
	}

	return &res, nil
}

func (g *userGRPCServer) ChangePassword(ctx context.Context, req *gen.ChangePasswordRequest) (*gen.ChangePasswordResponse, error) {
	slog.Info("PUT: ChangePassword")
	if len(req.NewPassword) < 8 {
		return nil, errors.New("Password must be at least 8 characters")
	}

	tokenString, err := sharedkernel.GetToken(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.GetToken")
	}

	tokenClaims, err := sharedkernel.ParseToken(ctx, tokenString)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.ParseToken")
	}

	ID, err := uuid.Parse(tokenClaims.Id)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse")
	}

	err = g.uc.ChangePassword(ctx, &domain.ChangePasswordModel{
		ID:              ID,
		OldPasswordHash: req.GetOldPassword(),
		NewPasswordHash: req.GetNewPassword(),
	})
	if err != nil {
		return nil, errors.Wrap(err, "usecase.ChangePassword")
	}

	return &gen.ChangePasswordResponse{}, nil
}

func (g *userGRPCServer) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error) {
  slog.Info("POST: RefreshToken")
	tokenString, err := sharedkernel.GetToken(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "sharedkernel.GetToken")
	}

	tokens, err := g.uc.RefreshToken(ctx, &domain.Token{
		AccessToken:  tokenString,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return nil, errors.Wrap(err, "uc.RefreshToken")
	}

	return &gen.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
