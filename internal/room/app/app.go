package app

import (
	"github.com/Watch2Gather/server/cmd/room/config"
	roomsUC "github.com/Watch2Gather/server/internal/room/usecases/rooms"
	"github.com/Watch2Gather/server/pkg/postgres"
	"github.com/Watch2Gather/server/proto/gen"
)

type App struct {
	Cfg             *config.Config
	PG              postgres.DBEngine
	UC              roomsUC.UseCase
	roomsGRPCServer gen.RoomServiceServer
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc roomsUC.UseCase,
	roomsGRPCServer gen.RoomServiceServer,
) *App {
	return &App{
		Cfg:             cfg,
		PG:              pg,
		UC:              uc,
		roomsGRPCServer: roomsGRPCServer,
	}
}
