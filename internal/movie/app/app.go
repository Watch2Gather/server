package app

import (
	"github.com/Watch2Gather/server/cmd/movie/config"
	moviesUC "github.com/Watch2Gather/server/internal/movie/usecases/movies"
	"github.com/Watch2Gather/server/pkg/postgres"
	"github.com/Watch2Gather/server/proto/gen"
)

type App struct {
	Cfg              *config.Config
	PG               postgres.DBEngine
	UC               moviesUC.UseCase
	moviesGRPCServer gen.MovieServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc moviesUC.UseCase,
	moviesGRPCServer gen.MovieServiceServer,
) *App {
	return &App{
		Cfg:              cfg,
		PG:               pg,
		UC:               uc,
		moviesGRPCServer: moviesGRPCServer,
	}
}
