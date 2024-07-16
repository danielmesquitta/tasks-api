-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = ?
LIMIT 1;
-- name: CreateUser :execresult
INSERT INTO users (role, name, email, password)
VALUES (?, ?, ?, ?);