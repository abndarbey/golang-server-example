-- name: ProfileInsert :one
INSERT INTO profiles (
  user_id,
  referral_code,
  wallet_points
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: ProfileGetByUser :one
SELECT * FROM profiles
WHERE user_id = $1 LIMIT 1;

-- name: UserGetByReferralCode :one
SELECT * FROM profiles
WHERE referral_code = $1 LIMIT 1;