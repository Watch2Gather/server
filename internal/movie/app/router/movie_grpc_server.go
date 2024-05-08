package router

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Watch2Gather/server/cmd/movie/config"
	"github.com/Watch2Gather/server/internal/movie/usecases/movies"
	"github.com/Watch2Gather/server/proto/gen"
)

type movieGRPCServer struct {
	cfg *config.Config
	uc  movies.UseCase
}

var _ gen.MovieServiceServer = (*movieGRPCServer)(nil)

var MovieGRPCServerSet = wire.NewSet(NewGRPCMovieServer)

func NewGRPCMovieServer(
	grpcServer *grpc.Server,
	cfg *config.Config,
	uc movies.UseCase,
) gen.MovieServiceServer {
	svc := movieGRPCServer{
		cfg: cfg,
		uc:  uc,
	}

	gen.RegisterMovieServiceServer(grpcServer, &svc)

	reflection.Register(grpcServer)

	return &svc
}

func (g *movieGRPCServer) GetAllMovies(ctx context.Context, req *gen.GetAllMoviesRequest) (*gen.GetAllMoviesResponse, error) {
	slog.Info("GET: GetAllMovies")

	// movies, err := g.uc.GetAllMovies(ctx)
	// if err != nil {
	// 	slog.Error("Caught error", "trace", errors.Wrap(err, "uc.GetAllMovies"))
	// 	return nil, sharedkernel.ErrServer
	// }
	//
	// _ = movies

	res := &gen.GetAllMoviesResponse{
		Movie: []*gen.ShortMovie{
			{
				Title:      "Lord of the Rings 1",
				Id:         uuid.New().String(),
				KpRating:   45,
				KpId:       123,
				PosterPath: "lord_of_the_rings_1",
			},
			{
				Title:      "Lord of the Rings 2",
				Id:         uuid.New().String(),
				KpRating:   46,
				KpId:       124,
				PosterPath: "lord_of_the_rings_2",
			},
			{
				Title:      "Lord of the Rings 3",
				Id:         uuid.New().String(),
				KpRating:   47,
				KpId:       125,
				PosterPath: "lord_of_the_rings_3",
			},
			{
				Title:      "Shrek 1",
				Id:         uuid.New().String(),
				KpRating:   47,
				KpId:       126,
				PosterPath: "shrek_1",
			},
			{
				Title:      "Shrek 2",
				Id:         uuid.New().String(),
				KpRating:   37,
				KpId:       127,
				PosterPath: "shrek_2",
			},
			{
				Title:      "Shrek 3",
				Id:         uuid.New().String(),
				KpRating:   47,
				KpId:       128,
				PosterPath: "shrek_3",
			},
		},
	}

	return res, nil
}

func (g *movieGRPCServer) GetMovie(ctx context.Context, req *gen.GetMovieRequest) (*gen.Movie, error) {
	slog.Info("GET: GetMovie")

	// id, err := uuid.Parse(req.GetId())
	// if err != nil {
	// 	slog.Error("Caught error", "trace", errors.Wrap(err, "uuid.Parse"))
	// 	return nil, sharedkernel.ErrServer
	// }
	//
	// movies, err := g.uc.GetMovieInfo(ctx, id)
	// if err != nil {
	// 	slog.Error("Caught error", "trace", errors.Wrap(err, "uc.GetMovieInfo"))
	// 	return nil, sharedkernel.ErrServer
	// }
	//
	// _ = movies

	res := &gen.Movie{
		Id:          uuid.New().String(),
		Title:       "Shrek 1",
		Description: "Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		KpRating:    47,
		ImdbRating:  98,
		KpId:        126,
		Year:        2002,
		PosterPath:  "shrek_1",
		Country:     "USA",
		ReviewCount: 1000,
	}

	return res, nil
}

func (g *movieGRPCServer) GetMoviePoster(_ context.Context, _ *gen.GetMoviePosterRequest) (_ *gen.GetMoviePosterResponse, _ error) {
	panic("not implemented") // TODO: Implement
}
