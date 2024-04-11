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

const addParticipant = `-- name: AddParticipant :exec
INSERT INTO "app".participants (
  room_id, user_id
) VALUES (
  $1, $2
)
`

type AddParticipantParams struct {
	RoomID uuid.UUID `json:"room_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) AddParticipant(ctx context.Context, arg AddParticipantParams) error {
	_, err := q.db.ExecContext(ctx, addParticipant, arg.RoomID, arg.UserID)
	return err
}

const createMessage = `-- name: CreateMessage :one
INSERT INTO "app".messages (
  room_id, user_id, content
) VALUES (
  $1, $2, $3
)
RETURNING id, room_id, user_id, content, created_at
`

type CreateMessageParams struct {
	RoomID  uuid.UUID `json:"room_id"`
	UserID  uuid.UUID `json:"user_id"`
	Content string    `json:"content"`
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (AppMessage, error) {
	row := q.db.QueryRowContext(ctx, createMessage, arg.RoomID, arg.UserID, arg.Content)
	var i AppMessage
	err := row.Scan(
		&i.ID,
		&i.RoomID,
		&i.UserID,
		&i.Content,
		&i.CreatedAt,
	)
	return i, err
}

const createRoom = `-- name: CreateRoom :one
INSERT INTO "app".rooms (
  name, owner_id
) VALUES (
  $1, $2
)
RETURNING id, name, owner_id, movie_id, timecode
`

type CreateRoomParams struct {
	Name    string    `json:"name"`
	OwnerID uuid.UUID `json:"owner_id"`
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) (AppRoom, error) {
	row := q.db.QueryRowContext(ctx, createRoom, arg.Name, arg.OwnerID)
	var i AppRoom
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.OwnerID,
		&i.MovieID,
		&i.Timecode,
	)
	return i, err
}

const deleteRoom = `-- name: DeleteRoom :exec
DELETE FROM "app".rooms
  WHERE id=$1 and owner_id=$2
`

type DeleteRoomParams struct {
	ID      uuid.UUID `json:"id"`
	OwnerID uuid.UUID `json:"owner_id"`
}

func (q *Queries) DeleteRoom(ctx context.Context, arg DeleteRoomParams) error {
	_, err := q.db.ExecContext(ctx, deleteRoom, arg.ID, arg.OwnerID)
	return err
}

const getMessagesByRoomId = `-- name: GetMessagesByRoomId :many
SELECT id, room_id, user_id, content, created_at FROM "app".messages
  WHERE room_id = $1
  ORDER BY created_at DESC
  LIMIT $2
  OFFSET $3
`

type GetMessagesByRoomIdParams struct {
	RoomID uuid.UUID `json:"room_id"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (q *Queries) GetMessagesByRoomId(ctx context.Context, arg GetMessagesByRoomIdParams) ([]AppMessage, error) {
	rows, err := q.db.QueryContext(ctx, getMessagesByRoomId, arg.RoomID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppMessage
	for rows.Next() {
		var i AppMessage
		if err := rows.Scan(
			&i.ID,
			&i.RoomID,
			&i.UserID,
			&i.Content,
			&i.CreatedAt,
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

const getParticipantsByRoomId = `-- name: GetParticipantsByRoomId :many
SELECT username, avatar, id FROM "app".users as u
INNER JOIN "app".participants as p
  on u.id = p.user_id
  WHERE p.room_id=$1
`

type GetParticipantsByRoomIdRow struct {
	Username string         `json:"username"`
	Avatar   sql.NullString `json:"avatar"`
	ID       uuid.UUID      `json:"id"`
}

func (q *Queries) GetParticipantsByRoomId(ctx context.Context, roomID uuid.UUID) ([]GetParticipantsByRoomIdRow, error) {
	rows, err := q.db.QueryContext(ctx, getParticipantsByRoomId, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetParticipantsByRoomIdRow
	for rows.Next() {
		var i GetParticipantsByRoomIdRow
		if err := rows.Scan(&i.Username, &i.Avatar, &i.ID); err != nil {
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

const getRoomsByUserId = `-- name: GetRoomsByUserId :many
SELECT id, name, owner_id, movie_id, timecode FROM "app".rooms as r
INNER JOIN "app".participants as p
  ON r.id = p.room_id
  WHERE p.user_id = $1
`

func (q *Queries) GetRoomsByUserId(ctx context.Context, userID uuid.UUID) ([]AppRoom, error) {
	rows, err := q.db.QueryContext(ctx, getRoomsByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []AppRoom
	for rows.Next() {
		var i AppRoom
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.OwnerID,
			&i.MovieID,
			&i.Timecode,
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

const removeParticipant = `-- name: RemoveParticipant :exec
DELETE FROM "app".participants
  WHERE room_id=$1 and user_id=$2
`

type RemoveParticipantParams struct {
	RoomID uuid.UUID `json:"room_id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) RemoveParticipant(ctx context.Context, arg RemoveParticipantParams) error {
	_, err := q.db.ExecContext(ctx, removeParticipant, arg.RoomID, arg.UserID)
	return err
}

const updateRoom = `-- name: UpdateRoom :exec
UPDATE "app".rooms
  SET
    name = COALESCE($1, name),
    movie_id = COALESCE($2, movie_id)
  WHERE id = $3
`

type UpdateRoomParams struct {
	Name    sql.NullString `json:"name"`
	MovieID uuid.NullUUID  `json:"movie_id"`
	ID      uuid.NullUUID  `json:"id"`
}

func (q *Queries) UpdateRoom(ctx context.Context, arg UpdateRoomParams) error {
	_, err := q.db.ExecContext(ctx, updateRoom, arg.Name, arg.MovieID, arg.ID)
	return err
}