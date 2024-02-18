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

const checkPassword = `-- name: CheckPassword :one
SELECT EXISTS (
  SELECT id FROM "app".users
    WHERE id=$1 AND pwd_hash=$2
)
`

type CheckPasswordParams struct {
	ID      uuid.UUID `json:"id"`
	PwdHash string    `json:"pwd_hash"`
}

func (q *Queries) CheckPassword(ctx context.Context, arg CheckPasswordParams) (bool, error) {
	row := q.db.QueryRowContext(ctx, checkPassword, arg.ID, arg.PwdHash)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createUser = `-- name: CreateUser :exec
INSERT INTO "app".users (
  id, username, email, pwd_hash
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, username, email, pwd_hash, avatar
`

type CreateUserParams struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	PwdHash  string    `json:"pwd_hash"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Email,
		arg.PwdHash,
	)
	return err
}

const loginUser = `-- name: LoginUser :one
SELECT id, username, email, avatar FROM "app".users
  WHERE username=$1 AND pwd_hash=$2
`

type LoginUserParams struct {
	Username string `json:"username"`
	PwdHash  string `json:"pwd_hash"`
}

type LoginUserRow struct {
	ID       uuid.UUID      `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	Avatar   sql.NullString `json:"avatar"`
}

func (q *Queries) LoginUser(ctx context.Context, arg LoginUserParams) (LoginUserRow, error) {
	row := q.db.QueryRowContext(ctx, loginUser, arg.Username, arg.PwdHash)
	var i LoginUserRow
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Email,
		&i.Avatar,
	)
	return i, err
}
