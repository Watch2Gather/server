-- name: CreateUser :one
INSERT INTO "app".users (
  username, email, pwd_hash
) VALUES (
  $1, $2, $3
)
RETURNING username, email, avatar;

-- name: GetUserByUsername :one
SELECT id, username, email, pwd_hash, avatar FROM "app".users
  WHERE username=$1;

-- name: GetUserByID :one
SELECT id, username, email, avatar, token FROM "app".users
  WHERE id=$1;

-- name: GetPasswordById :one
SELECT pwd_hash FROM "app".users
  WHERE id=$1;

-- name: UpdateUserPassword :exec
UPDATE "app".users
  SET pwd_hash = $2
  WHERE id = $1;

-- name: UpdateUser :one
UPDATE "app".users
  SET
    username = COALESCE(sqlc.narg(username), username),
    email = COALESCE(sqlc.narg(email), email),
    avatar = COALESCE(sqlc.narg(avatar), avatar)
  WHERE id = sqlc.arg(id)
RETURNING username, email, avatar;

-- name: UpdateToken :exec
UPDATE "app".users
  SET token = $2
  WHERE id = $1;

-- name: GetUserTokenById :one
SELECT token FROM "app".users
  WHERE id = $1;

-- name: GetFriendList :many
SELECT user_id_2 as id, u.username, u.avatar FROM "app".friends as f
  join "app".users as u on f.user_id_2 = u.id
  WHERE f.user_id_1 = $1;

-- name: AddFriendById :exec
INSERT INTO "app".friends (
  user_id_1, user_id_2
) VALUES ( $1, $2 );
