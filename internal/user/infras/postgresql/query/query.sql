-- name: CreateUser :one
INSERT INTO "app".users (
  id, username, email, pwd_hash
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: LoginUser :one
SELECT id, username, email, avatar FROM "app".users
  WHERE username=$1 AND pwd_hash=$2;

-- name: CheckPassword :one
SELECT EXISTS (
  SELECT id FROM "app".users
    WHERE id=$1 AND pwd_hash=$2
);
