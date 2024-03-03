package app

import (
	"github.com/Watch2Gather/server/cmd/user/config"
	"github.com/Watch2Gather/server/pkg/postgres"
)

type App struct {
	Cfg *config.Config
	PG  postgres.DBEngine

	// TODO Add other services here

	//
}
