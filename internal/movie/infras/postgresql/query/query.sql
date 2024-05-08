-- name: GetAllMovies :many
SELECT id, title, kp_rating, kp_id, poster_path FROM "app".movies;

-- name: GetMovieByID :one
SELECT id, title, kp_rating, kp_id, poster_path  FROM "app".movies
WHERE id = $1;

-- name: GetPosterPath :one
SELECT poster_path FROM "app".movies
WHERE id = $1;
