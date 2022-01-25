-- name: UserInsert :one
INSERT INTO users (
  first_name,
  last_name,
  email,
  phone,
  password_hash
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: UserList :many
SELECT * FROM users
ORDER BY id;

-- name: UserGetByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: UserGetByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: UserGetByPhone :one
SELECT * FROM users
WHERE phone = $1 LIMIT 1;

-- name: UserGetProfile :one
SELECT * FROM users
RIGHT JOIN profiles
ON users.id = profiles.user_id
WHERE id = $1
LIMIT 1;

-- name: UserGetWithReferrals :one
SELECT * FROM users
RIGHT JOIN referrals
ON users.id = referrals.referee_id
WHERE users.id = $1
LIMIT 1;

-- name: UserDelete :exec
DELETE FROM users
WHERE id = $1;