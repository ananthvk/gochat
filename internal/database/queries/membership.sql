-- Check if user is part of a group

-- name: CheckMembership :one
SELECT EXISTS(
    SELECT 1 FROM grp_membership
    WHERE
        grp_id = sqlc.arg('grp_id')
            AND
        usr_id = sqlc.arg('usr_id')
    )
;

-- Add a user to a group  

-- name: CreateMembership :one
INSERT INTO grp_membership (grp_id, usr_id, role)
VALUES (sqlc.arg('grp_id'), sqlc.arg('usr_id'), 'member')
RETURNING joined_at;

-- Remove a user from a group

-- name: DeleteMembership :exec
DELETE FROM grp_membership
WHERE
    grp_id = sqlc.arg('grp_id')
        AND
    usr_id = sqlc.arg('usr_id')
;

-- Get all users in a group

-- name: GetGroupMemberships :many
SELECT * FROM grp_membership
WHERE
grp_id = sqlc.arg('grp_id')
;


-- Get all groups the user is a part of

-- name: GetUserMemberships :many
SELECT * FROM grp_membership
WHERE
usr_id = sqlc.arg('usr_id')
;