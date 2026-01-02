-- name: GetGroup :one
SELECT * FROM grp
WHERE id = sqlc.arg('id') LIMIT 1;

-- Returns detailed information about all groups the user is part of
-- Also sorts the returned groups by the last message sent timestamp
-- Also returns the last message (if any) of the group
-- Also joins with the user table, and returns the sender name, so that it can be directly rendered in the sidebar
-- name: GetGroups :many
SELECT 
    g.*,
    mem.role,
    mem.joined_at,
    m.content AS last_message_content,
    m.created_at AS last_message_created_at,
    last_message_id,
    m.sender_id AS last_message_sender_id,
    m.type AS last_message_type,
    u.name AS last_message_sender_name
FROM grp AS g
INNER JOIN grp_membership AS mem
    ON g.id = mem.grp_id
LEFT JOIN (
    SELECT DISTINCT ON (grp_id)
        grp_id,
        id AS last_message_id
    FROM message
    ORDER BY grp_id, id DESC
) AS lm 
    ON lm.grp_id = g.id
LEFT JOIN message as m
    ON m.id = lm.last_message_id
LEFT JOIN usr AS u
    ON u.id = m.sender_id
WHERE mem.usr_id  = sqlc.arg('usr_id')
ORDER BY lm.last_message_id DESC NULLS LAST;

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