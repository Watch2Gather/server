-- name: CreateRoom :one
INSERT INTO "app".rooms (
  name, owner_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetRoomsByUserId :many
SELECT id, name, owner_id, movie_id, timecode, (
  SELECT COUNT(*) FROM "app".participants as p
  WHERE p.room_id = r.id
)
FROM "app".rooms as r
INNER JOIN "app".participants AS p
  ON r.id = p.room_id
  WHERE p.user_id = $1;

-- name: GetParticipantsByRoomId :many
SELECT username, avatar, id FROM "app".users AS u
INNER JOIN "app".participants AS p
  ON u.id = p.user_id
  WHERE p.room_id=$1;

-- name: DeleteRoom :exec
DELETE FROM "app".rooms
  WHERE id=$1 AND owner_id=$2;

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
SELECT m.id AS m_id, m.content, m.created_at, u.id AS u_id, u.username, u.avatar  FROM "app".messages AS m
  JOIN "app".users AS u
    ON m.user_id = u.id
   	WHERE m.room_id = $1
   	ORDER BY m.created_at ASC;

-- name: AddParticipant :exec
INSERT INTO "app".participants (
  room_id, user_id
) VALUES (
  $1, $2
);

-- name: GetRoomOwner :one
SELECT owner_id FROM "app".rooms
WHERE id=$1;

-- name: RemoveParticipant :exec
DELETE FROM "app".participants
  WHERE room_id=$1 and user_id=$2;
