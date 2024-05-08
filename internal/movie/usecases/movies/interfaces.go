package movies

import (
	"context"

	"github.com/google/uuid"

	"github.com/Watch2Gather/server/internal/movie/domain"
)

type (
	MovieRepo interface {
	GetAllMovies(context.Context) ([]*domain.ShortMovieModel, error)
	GetMovieInfo(context.Context, uuid.UUID) (*domain.MovieModel, error)
	GetMoviePosterPath(context.Context, string) (*domain.Poster, error)
	GetMoviePosterPicture(context.Context, string)
	}
	UseCase interface {
	GetAllMovies(context.Context) ([]*domain.ShortMovieModel, error)
	GetMovieInfo(context.Context, uuid.UUID) (*domain.MovieModel, error)
	GetMoviePoster(context.Context, string) (*domain.Poster, error)
	}
)
