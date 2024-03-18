package app

import (
	"google.golang.org/grpc"

	"github.com/google/wire"

	"github.com/Watch2Gather/server/cmd/user/config"
	"github.com/Watch2Gather/server/internal/user/infras/repo"
	usersUC "github.com/Watch2Gather/server/internal/user/usecases/users"
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
		// TODO: add router
		repo.RepositorySet,
		usersUC.UseCaseSet,
	))
}

func dbEngineFunc(url postgres.DBConnString) (postgres.DBEngine, func(), error) {
	db, err := postgres.NewPostgresDB(url)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
