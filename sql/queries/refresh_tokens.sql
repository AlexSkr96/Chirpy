-- name: SaveRefreshToken :exec
INSERT INTO refresh_tokens (user_id, token)
VALUES ($1, $2);

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :exec
update refresh_tokens set revoked_at = now() where token = $1;
