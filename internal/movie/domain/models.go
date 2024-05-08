package domain

import (
	"github.com/google/uuid"
)

type ShortMovieModel struct {
	Title      string
	PosterPath string
	KpRating   int
	KpID       int
	ID         uuid.UUID
}

type MovieModel struct {
	Desription  string
	Country     string
	Info        ShortMovieModel
	Year        int
	ReviewCount int
}

type Poster []byte
