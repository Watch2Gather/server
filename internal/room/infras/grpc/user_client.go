package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/Watch2Gather/server/cmd/room/config"
	"github.com/Watch2Gather/server/internal/room/domain"
	"github.com/Watch2Gather/server/proto/gen"
)

type userInfoGRPCClient struct {
	conn *grpc.ClientConn
}

var _ domain.UserInfoDomainService = (*userInfoGRPCClient)(nil)

var UserInfoGRPCClientSet = wire.NewSet(NewGRPCUserInfoClient)

func NewGRPCUserInfoClient(cfg *config.Config) (domain.UserInfoDomainService, error) {
	conn, err := grpc.Dial(cfg.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &userInfoGRPCClient{
		conn: conn,
	}, nil
}

func (u *userInfoGRPCClient) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.UserModel, error) {
	c := gen.NewUserInfoServiceClient(u.conn)
	res, err := c.GetUserInfo(ctx, &gen.GetUserInfoRequest{Id: id.String()})
	if err != nil {
		return nil, errors.Wrap(err, "userInfo.GetUserInfo")
	}

	uid, err := uuid.Parse(res.User.Id)
	if err != nil {
		return nil, errors.Wrap(err, "uuid.Parse")
	}

	return &domain.UserModel{
		ID:     uid,
		Name:   res.User.Username,
		Avatar: res.User.Avatar,
	}, nil
}
