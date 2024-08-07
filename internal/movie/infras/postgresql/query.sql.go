// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package postgresql

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const getAllMovies = `-- name: GetAllMovies :many
SELECT id, title, kp_rating, kp_id, poster_path FROM "app".movies
`

type GetAllMoviesRow struct {
	ID         uuid.UUID      `json:"id"`
	Title      string         `json:"title"`
	KpRating   sql.NullString `json:"kp_rating"`
	KpID       sql.NullInt32  `json:"kp_id"`
	PosterPath string         `json:"poster_path"`
}

func (q *Queries) GetAllMovies(ctx context.Context) ([]GetAllMoviesRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllMovies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllMoviesRow
	for rows.Next() {
		var i GetAllMoviesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.KpRating,
			&i.KpID,
			&i.PosterPath,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMovieByID = `-- name: GetMovieByID :one
SELECT id, title, kp_rating, kp_id, poster_path  FROM "app".movies
WHERE id = $1
`

type GetMovieByIDRow struct {
	ID         uuid.UUID      `json:"id"`
	Title      string         `json:"title"`
	KpRating   sql.NullString `json:"kp_rating"`
	KpID       sql.NullInt32  `json:"kp_id"`
	PosterPath string         `json:"poster_path"`
}

func (q *Queries) GetMovieByID(ctx context.Context, id uuid.UUID) (GetMovieByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getMovieByID, id)
	var i GetMovieByIDRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.KpRating,
		&i.KpID,
		&i.PosterPath,
	)
	return i, err
}

const getPosterPath = `-- name: GetPosterPath :one
SELECT poster_path FROM "app".movies
WHERE id = $1
`

func (q *Queries) GetPosterPath(ctx context.Context, id uuid.UUID) (string, error) {
	row := q.db.QueryRowContext(ctx, getPosterPath, id)
	var poster_path string
	err := row.Scan(&poster_path)
	return poster_path, err
}
