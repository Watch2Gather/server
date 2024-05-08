//go:build wireinject

package app

import (
	"google.golang.org/grpc"

	"github.com/google/wire"

	"github.com/Watch2Gather/server/cmd/movie/config"
	"github.com/Watch2Gather/server/internal/movie/app/router"
	"github.com/Watch2Gather/server/internal/movie/infras/repo"
	moviesUC "github.com/Watch2Gather/server/internal/movie/usecases/movies"
	"github.com/Watch2Gather/server/pkg/postgres"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	grpcServer *grpc.Server,
) (*App, func(), error) {
	panic(wire.Build(
		New,
		dbEngineFunc,

		router.MovieGRPCServerSet,
		repo.RepositorySet,
		moviesUC.UseCaseSet,
	))
}

func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
