package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/pkg/errors"

	"github.com/Watch2Gather/server/internal/movie/domain"
	"github.com/Watch2Gather/server/internal/movie/infras/postgresql"
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

func (m *movieRepo) GetAllMovies(ctx context.Context) ([]*domain.ShortMovieModel, error) {
	querier := postgresql.New(m.pg.GetDB())

	movies, err := querier.GetAllMovies(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "querier.GetAllMovies")
	}

	shortMovies := []*domain.ShortMovieModel{}

	for _, movie := range movies {
		shortMovies = append(shortMovies, &domain.ShortMovieModel{
			Title:      movie.Title,
			PosterPath: movie.PosterPath,
			// TODO change data type in db
			KpRating: 70,
			KpID:     int(movie.KpID.Int32),
			ID:       movie.ID,
		})
	}

	return shortMovies, nil
}

func (movieRepo) GetMovieInfo(_ context.Context, _ uuid.UUID) (_ *domain.MovieModel, _ error) {
	panic("not implemented") // TODO: Implement
}

func (movieRepo) GetMoviePosterPath(_ context.Context, _ string) (string, error) {
	panic("not implemented") // TODO: Implement
}
