package app

import (
	"github.com/Watch2Gather/server/cmd/user/config"
	usersUC "github.com/Watch2Gather/server/internal/user/usecases/users"
	"github.com/Watch2Gather/server/pkg/postgres"
)

type App struct {
	Cfg *config.Config
	PG  postgres.DBEngine

	UC usersUC.UseCase
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,

	uc usersUC.UseCase,
) *App {
	return &App{
		Cfg: cfg,
		PG:  pg,
		UC:  uc,
	}
}
