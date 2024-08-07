package router

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/Watch2Gather/server/cmd/movie/config"
	"github.com/Watch2Gather/server/internal/movie/usecases/movies"
	sharedkernel "github.com/Watch2Gather/server/internal/pkg/shared_kernel"
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

	movies, err := g.uc.GetAllMovies(ctx)
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.GetAllMovies"))
		return nil, sharedkernel.ErrServer
	}

	res := &gen.GetAllMoviesResponse{}

	for _, movie := range movies {
		res.Movies = append(res.Movies, &gen.ShortMovie{
			Id:         movie.ID.String(),
			Title:      movie.Title,
			KpRating:   int32(movie.KpRating),
			KpId:       int32(movie.KpID),
			PosterPath: movie.PosterPath,
		})
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

func (g *movieGRPCServer) GetMoviePoster(ctx context.Context, req *gen.GetMoviePosterRequest) (*gen.GetMoviePosterResponse, error) {
	slog.Info("GET: GetMoviePoster")

	poster, err := g.uc.GetMoviePoster(ctx, req.GetFilePath())
	if err != nil {
		slog.Error("Caught error", "trace", errors.Wrap(err, "uc.GetMoviePoster"))
		return nil, nil
	}

	res := &gen.GetMoviePosterResponse{
		Poster: *poster,
	}

	return res, nil
}
