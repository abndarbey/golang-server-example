-- name: SessionCreate :one
INSERT INTO sessions (
  user_id,
  token,
  expires,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: SessionGetByToken :one
SELECT * FROM sessions
WHERE token = $1 LIMIT 1;