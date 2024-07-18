-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;
-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = ?
LIMIT 1;
-- name: CreateUser :exec
INSERT INTO users (role, name, email, password)
VALUES (?, ?, ?, ?);