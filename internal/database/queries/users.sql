-- name: GetUserByEmail :one
SELECT * FROM usr
WHERE email = sqlc.arg('email') LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM usr
WHERE username = sqlc.arg('username') LIMIT 1;

-- name: GetUserById :one
SELECT * FROM usr
where id = sqlc.arg('id') limit 1;

-- name: CreateUser :one
INSERT INTO usr (id, name, username, email, password, activated)
VALUES (
    sqlc.arg('id'),
    sqlc.arg('name'),
    sqlc.arg('username'),
    sqlc.arg('email'),
    sqlc.arg('password'),
    sqlc.arg('activated')
)
RETURNING id, created_at;

-- name: UpdateUser :one
UPDATE usr
SET
    name = coalesce(sqlc.narg('name'), name),
    username = coalesce(sqlc.narg('username'), username),
    email = coalesce(sqlc.narg('email'), email),
    password = coalesce(sqlc.narg('password'), password),
    activated = coalesce(sqlc.narg('activated'), activated)
WHERE id = sqlc.arg('id')
RETURNING *;