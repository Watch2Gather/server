package app

import (
	"github.com/Watch2Gather/server/cmd/room/config"
	"github.com/Watch2Gather/server/internal/room/domain"
	roomsUC "github.com/Watch2Gather/server/internal/room/usecases/rooms"
	"github.com/Watch2Gather/server/pkg/postgres"
	"github.com/Watch2Gather/server/proto/gen"
)

type App struct {
	Cfg               *config.Config
	PG                postgres.DBEngine
	UC                roomsUC.UseCase
	RoomsGRPCServer   gen.RoomServiceServer
	UserInfoDomainSvc domain.UserInfoDomainService
}

func New(
	cfg *config.Config,
	pg postgres.DBEngine,
	uc roomsUC.UseCase,
	roomsGRPCServer gen.RoomServiceServer,
	userInfoDomainSvc domain.UserInfoDomainService,
) *App {
	return &App{
		Cfg:               cfg,
		PG:                pg,
		UC:                uc,
		RoomsGRPCServer:   roomsGRPCServer,
		UserInfoDomainSvc: userInfoDomainSvc,
	}
}
