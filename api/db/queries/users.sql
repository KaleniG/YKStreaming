-- name: GetUserDataByRememberToken :one
SELECT id, name, email 
FROM users 
WHERE remember_token = sqlc.arg(remember_token)
LIMIT 1;

-- name: GetUserDataById :one
SELECT name, email
FROM users
WHERE id = sqlc.arg(user_id) 
LIMIT 1;

-- name: GetUserCredentialsByEmail :one
SELECT id, password_hash
FROM users
WHERE email = sqlc.arg(email) 
LIMIT 1;

-- name: GetUserIdByRememberToken :one
SELECT 1 
FROM users 
WHERE remember_token = sqlc.arg(remember_token)
LIMIT 1;

-- name: UpdateUserRememberToken :exec
UPDATE users 
SET remember_token = sqlc.arg(remember_token) 
WHERE id = sqlc.arg(user_id);

-- name: AddUser :one
INSERT INTO users (name, email, password_hash) 
VALUES 
(
  sqlc.arg(name), 
  sqlc.arg(email),
  sqlc.arg(password_hash)
) 
ON CONFLICT (email) DO NOTHING 
RETURNING id;

-- name: CheckUserById :one
SELECT 1
FROM users
WHERE id = sqlc.arg(user_id) 
LIMIT 1;

-- name: ResetRememberToken :exec
UPDATE users 
SET remember_token = NULL 
WHERE remember_token = sqlc.arg(remember_token);

-- name: ResetUserRememberToken :exec
UPDATE users 
SET remember_token = NULL 
WHERE id = sqlc.arg(user_id);

-- name: RemoveUserByEmail :exec
DELETE FROM users
WHERE email = sqlc.arg(email);