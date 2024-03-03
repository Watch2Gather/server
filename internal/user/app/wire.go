package app

import (
	"google.golang.org/grpc"

	"github.com/Watch2Gather/server/cmd/user/config"
	"github.com/Watch2Gather/server/pkg/postgres"
)

func InitApp(
	cfg *config.Config,
	dbConnStr postgres.DBConnString,
	grpcServer *grpc.Server,
) (*App, func(), error)
