-- name: GetGroup :one
SELECT * FROM grp
WHERE id = sqlc.arg('id') LIMIT 1;

-- name: GetGroups :many
SELECT * FROM grp;

-- name: CreateGroup :one
INSERT INTO grp (id, name, description)
VALUES (sqlc.narg('id'), sqlc.arg('name'), sqlc.arg('description'))
RETURNING id;

-- name: DeleteGroup :exec
DELETE FROM grp
WHERE id = sqlc.arg('id');

-- name: UpdateGroupById :one
UPDATE grp 
SET
    name = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: CheckGroupExists :one
SELECT EXISTS(SELECT 1 FROM grp WHERE id = sqlc.arg('id')) as exists;