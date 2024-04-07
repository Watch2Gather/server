package router

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

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

	if len(req.GetUsername()) < 3 || len(req.GetUsername()) > 36 {
		return nil, status.Error(codes.InvalidArgument, "Username must be between 3 and 36 characters")
	}

	if len(req.GetEmail()) < 3 || len(req.GetEmail()) > 254 {
		return nil, status.Error(codes.InvalidArgument, "Email must be between 3 and 254 characters")
	}

	if len(req.GetPassword()) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters")
	}

	err := g.uc.Register(ctx, &domain.RegisterModel{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
		Email:    req.GetEmail(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.Register"))
		return nil, sharedkernel.ErrServer
	}

	res := gen.RegisterUserResponse{}

	return &res, nil
}

func (g *userGRPCServer) LoginUser(ctx context.Context, req *gen.LoginUserRequest) (*gen.LoginUserResponse, error) {
	slog.Info("POST: LoginUser")

	if len(req.GetUsername()) < 3 || len(req.GetUsername()) > 36 {
		return nil, status.Error(codes.InvalidArgument, "Username must be between 3 and 36 characters")
	}

	if len(req.GetPassword()) < 8 {
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters")
	}
	res := gen.LoginUserResponse{}
	token, err := g.uc.Login(ctx, &domain.LoginModel{
		Username: req.GetUsername(),
		Password: req.GetPassword(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			return nil, status.Error(codes.Unauthenticated, "Wrong username or password")
		}
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.Login"))
		return nil, sharedkernel.ErrServer
	}

	res.AccessToken = token.AccessToken
	res.RefreshToken = token.RefreshToken

	return &res, nil
}

func (g *userGRPCServer) ChangeUserData(ctx context.Context, req *gen.ChangeUserDataRequest) (*gen.ChangeUserDataResponse, error) {
	slog.Info("PUT: ChangeUserData")

	if len(req.GetUsername()) < 3 || len(req.GetUsername()) > 36 {
		return nil, status.Error(codes.InvalidArgument, "Username must be between 3 and 36 characters")
	}

	if len(req.GetEmail()) < 3 || len(req.GetEmail()) > 254 {
		return nil, status.Error(codes.InvalidArgument, "Email must be between 3 and 254 characters")
	}

	tokenString, err := sharedkernel.GetToken(ctx)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "sharedkernel.GetToken"))
		return nil, sharedkernel.ErrServer
	}

	tokenClaims, err := sharedkernel.ParseToken(ctx, tokenString)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "sharedkernel.ParseToken"))
		return nil, sharedkernel.ErrServer
	}

	id, err := uuid.Parse(tokenClaims.Id)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	user, err := g.uc.ChangeUserData(ctx, &domain.User{
		ID:       id,
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Avatar:   req.GetAvatar(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "User already exists")
		}
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.Login"))
		return nil, sharedkernel.ErrServer
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
		return nil, status.Error(codes.InvalidArgument, "Password must be at least 8 characters")
	}

	tokenString, err := sharedkernel.GetToken(ctx)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "sharedkernel.GetToken"))
		return nil, sharedkernel.ErrServer
	}

	tokenClaims, err := sharedkernel.ParseToken(ctx, tokenString)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "sharedkernel.ParseToken"))
		return nil, sharedkernel.ErrServer
	}

	ID, err := uuid.Parse(tokenClaims.Id)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
		return nil, sharedkernel.ErrServer
	}

	err = g.uc.ChangePassword(ctx, &domain.ChangePasswordModel{
		ID:              ID,
		OldPassword: req.GetOldPassword(),
		NewPassword: req.GetNewPassword(),
	})
	if err != nil {
		if errors.Is(err, domain.ErrInvalidPassword){
			return nil, status.Error(codes.InvalidArgument, "Invalid password")
		}
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.ChangePassword"))
		return nil, sharedkernel.ErrServer
	}

	return &gen.ChangePasswordResponse{}, nil
}

func (g *userGRPCServer) RefreshToken(ctx context.Context, req *gen.RefreshTokenRequest) (*gen.RefreshTokenResponse, error) {
	slog.Info("POST: RefreshToken")
	tokenString, err := sharedkernel.GetToken(ctx)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "sharedkernel.GetToken"))
		return nil, sharedkernel.ErrServer
	}

	tokens, err := g.uc.RefreshToken(ctx, &domain.Token{
		AccessToken:  tokenString,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.RefreshToken"))
		return nil, sharedkernel.ErrServer
	}

	return &gen.RefreshTokenResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
