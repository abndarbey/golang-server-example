-- name: AddressInsert :one
INSERT INTO addresses (
  user_id,
  line_1,
  line_2,
  line_3,
  city,
  state,
  country
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: AddressGetByUser :one
SELECT * FROM addresses
WHERE user_id = $1 LIMIT 1;