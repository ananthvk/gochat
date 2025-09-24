-- name: GetGroup :one
SELECT * FROM grp
WHERE id = $1 LIMIT 1;
