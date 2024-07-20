-- name: CreateUser :one
INSERT INTO auth (
  "username", 
  "password", 
  "profile_public_id", 
  "email", 
  "country", 
  "profile_picture", 
  "email_verification_token"
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM auth
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM auth
WHERE email = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM auth
WHERE username = $1 LIMIT 1;

-- name: GetUserByVerificationToken :one
SELECT * FROM auth
WHERE email_verification_token = $1 AND email_verified = FALSE LIMIT 1;

-- name: GetUserByPasswordToken :one
SELECT * FROM auth
WHERE password_reset_token = $1 AND password_reset_expires > NOW();

-- name: UpdateVerifyEmailField :exec
UPDATE auth 
SET email_verified = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdateEmailVerificationToken :execrows
UPDATE auth 
SET email_verification_token = $2, updated_at = NOW()
WHERE id = $1;

-- name: UpdatePasswordToken :exec
UPDATE auth 
SET 
  password_reset_token = $2, 
  password_reset_expires = $3,
  updated_at = NOW()
WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE auth 
SET password = $2, password_reset_token = NULL, updated_at = NOW()
WHERE id = $1;

