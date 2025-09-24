-- name: GetGroupByPublicId :one
SELECT * FROM grp
WHERE public_id = sqlc.arg('public_id') LIMIT 1;

-- name: GetGroups :many
SELECT * FROM grp;

-- name: CreateGroup :one
INSERT INTO grp (name, description, public_id)
VALUES (sqlc.arg('name'), sqlc.arg('description'), sqlc.narg('public_id'))
RETURNING id;

-- name: DeleteGroupByPublicId :exec
DELETE FROM grp
WHERE public_id = sqlc.arg('public_id');

-- name: UpdateGroupByPublicId :one
UPDATE grp 
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description)
WHERE public_id = sqlc.arg('public_id')
RETURNING *;