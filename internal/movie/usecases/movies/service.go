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

func (usecase) GetAllMovies(_ context.Context) (_ []*domain.ShortMovieModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (usecase) GetMovieInfo(_ context.Context, _ uuid.UUID) (_ *domain.MovieModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (usecase) GetMoviePoster(ctx context.Context, path string) (*[]byte, error) {
	// panic("not implemented") // TODO: Implement

	f, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}

	return &f, nil
}
