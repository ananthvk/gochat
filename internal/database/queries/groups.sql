-- name: GetGroup :one
SELECT * FROM grp
WHERE id = sqlc.arg('id') LIMIT 1;

-- Returns detailed information about all groups the user is part of
-- name: GetGroups :many
SELECT 
    g.*,
    mem.role,
    mem.joined_at 
FROM grp AS g
INNER JOIN grp_membership AS mem
    ON g.id = mem.grp_id
WHERE mem.usr_id  = sqlc.arg('usr_id');

-- name: CreateGroup :one
INSERT INTO grp (id, name, description, owner_id)
VALUES (sqlc.arg('id'), sqlc.arg('name'), sqlc.arg('description'), sqlc.arg('owner_id'))
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