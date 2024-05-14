package movies

import (
	"context"
	"os"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/Watch2Gather/server/internal/movie/domain"
)

type usecase struct {
	movieRepo MovieRepo
}

var _ UseCase = (*usecase)(nil)

var UseCaseSet = wire.NewSet(NewUseCase)

func NewUseCase(
	movieRepo MovieRepo,
) UseCase {
	return &usecase{
		movieRepo: movieRepo,
	}
}

var pathPrefix = os.Getenv("POSTER_PATH_PREFIX")

func (u *usecase) GetAllMovies(ctx context.Context) ([]*domain.ShortMovieModel, error) {
	movies, err := u.movieRepo.GetAllMovies(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "roomRepo.GetRoomsByUserID")
	}

	return movies, nil
}

func (usecase) GetMovieInfo(_ context.Context, _ uuid.UUID) (_ *domain.MovieModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (usecase) GetMoviePoster(ctx context.Context, path string) (*[]byte, error) {
	f, err := os.ReadFile(pathPrefix + path)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}

	return &f, nil
}
