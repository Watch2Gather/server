package router

import (
	"github.com/google/wire"

	"github.com/Watch2Gather/server/cmd/movie/config"
	"github.com/Watch2Gather/server/internal/movie/usecases/movies"

	_ "github.com/pion/webrtc/v4"
)

type movieWebRTCServer struct {
	cfg *config.Config
}

var MovieWebRTCServerSet = wire.NewSet(NewWebRTCMovieServer)

func NewWebRTCMovieServer(
	// grpcServer *grpc.Server,
	cfg *config.Config,
	uc movies.UseCase,
) {
}

func StartStreaming() {
}
