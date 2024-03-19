package app

import (
	"github.com/Watch2Gather/server/cmd/user/config"
	usersUC "github.com/Watch2Gather/server/internal/user/usecases/users"
	"github.com/Watch2Gather/server/pkg/postgres"
	"github.com/Watch2Gather/server/proto/gen"
)

type App struct {
	Cfg            *config.Config
	PG             postgres.DBEngine
	UC             usersUC.UseCase
	userGRPCServer gen.UserServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc usersUC.UseCase,
	userGRPCServer gen.UserServiceServer,
) *App {
	return &App{
		Cfg:            cfg,
		PG:             pg,
		UC:             uc,
		userGRPCServer: userGRPCServer,
	}
}
