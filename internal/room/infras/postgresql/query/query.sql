-- name: CreateRoom :one
INSERT INTO "app".rooms (
  name, owner_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetRoomsByUserId :many
SELECT id, name, owner_id, movie_id, timecode FROM "app".rooms as r
INNER JOIN "app".participants as p
  ON r.id = p.room_id
  WHERE p.user_id = $1;

-- name: GetParticipantsByRoomId :many
SELECT username, avatar, id FROM "app".users as u
INNER JOIN "app".participants as p
  on u.id = p.user_id
  WHERE p.room_id=$1;

-- name: DeleteRoom :exec
DELETE FROM "app".rooms
  WHERE id=$1 and owner_id=$2;

-- name: UpdateRoom :exec
UPDATE "app".rooms
  SET
    name = COALESCE(sqlc.narg(name), name),
    movie_id = COALESCE(sqlc.narg(movie_id), movie_id)
  WHERE id = sqlc.narg(id);

-- name: CreateMessage :one
INSERT INTO "app".messages (
  room_id, user_id, content
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetMessagesByRoomId :many
SELECT * FROM "app".messages
  WHERE room_id = $1
  ORDER BY created_at DESC
  LIMIT $2
  OFFSET $3;

-- name: AddParticipant :exec
INSERT INTO "app".participants (
  room_id, user_id
) VALUES (
  $1, $2
);

-- name: RemoveParticipant :exec
DELETE FROM "app".participants
  WHERE room_id=$1 and user_id=$2;
