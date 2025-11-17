-- name: GetMessage :one
SELECT * FROM message
WHERE
    id = sqlc.arg('id')
        AND
    grp_id = sqlc.arg('grp_id')
LIMIT 1;

-- name: GetMessagesInGroup :many
SELECT * FROM message
WHERE grp_id = sqlc.arg('grp_id');

-- name: DeleteMessage :exec
DELETE FROM message
WHERE 
    id = sqlc.arg('id')
        AND
    grp_id = sqlc.arg('grp_id')
;

-- name: CreateMessage :one
INSERT INTO message (id, type, grp_id, content)
VALUES (sqlc.arg('id'), sqlc.arg('type'), sqlc.arg('grp_id'), sqlc.arg('content'))
RETURNING id;