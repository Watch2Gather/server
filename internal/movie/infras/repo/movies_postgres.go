package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/google/wire"

	"github.com/Watch2Gather/server/internal/movie/domain"
	"github.com/Watch2Gather/server/internal/movie/usecases/movies"
	"github.com/Watch2Gather/server/pkg/postgres"
)

const _defaultEntityCap = 64

type movieRepo struct {
	pg postgres.DBEngine
}

var _ movies.MovieRepo = (*movieRepo)(nil)

var RepositorySet = wire.NewSet(NewMovieRepo)

func NewMovieRepo(
	pg postgres.DBEngine,
) movies.MovieRepo {
	return &movieRepo{pg: pg}
}

func (movieRepo) GetAllMovies(_ context.Context) (_ []*domain.ShortMovieModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (movieRepo) GetMovieInfo(_ context.Context, _ uuid.UUID) (_ *domain.MovieModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (movieRepo) GetMoviePosterPath(_ context.Context, _ string) (string, error) {
	panic("not implemented") // TODO: Implement
}
