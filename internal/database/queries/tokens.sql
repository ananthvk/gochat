-- name: CreateToken :exec
INSERT INTO token
    (hash, usr_id, expiry, scope)
VALUES (
        sqlc.arg('hash'),
        sqlc.arg('usr_id'),
        sqlc.arg('expiry'),
        sqlc.arg('scope')
);

-- name: VerifyToken :one
SELECT usr_id FROM token
WHERE
    scope = sqlc.arg('scope')
    AND hash = sqlc.arg('hash')
    AND expiry > NOW();

-- name: DeleteTokenByHash :exec
DELETE FROM token
WHERE hash = sqlc.arg('hash');

-- name: DeleteTokensByUserId :exec
DELETE FROM token
WHERE
    usr_id = sqlc.arg('usr_id')
        AND
    scope = sqlc.arg('scope');
