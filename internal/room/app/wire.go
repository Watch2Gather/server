//go:build wireinject
// +build wireinject

package app

import (
	"google.golang.org/grpc"

	"github.com/google/wire"

	"github.com/Watch2Gather/server/cmd/room/config"
	"github.com/Watch2Gather/server/internal/room/app/router"
	"github.com/Watch2Gather/server/internal/room/infras/repo"
	usersUC "github.com/Watch2Gather/server/internal/room/usecases/rooms"
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

		router.RoomGRPCServerSet,
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
