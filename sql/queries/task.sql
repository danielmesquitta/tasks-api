-- name: GetTaskByID :one
SELECT *
FROM tasks
WHERE id = ?
LIMIT 1;
-- name: ListTasks :many
SELECT *
FROM tasks;
-- name: ListTasksWithFilters :many
SELECT *
FROM tasks
WHERE assigned_to_user_id = ?;
-- name: CreateTask :exec
INSERT INTO tasks (summary, created_by_user_id, assigned_to_user_id)
VALUES (?, ?, ?);
-- name: UpdateTask :exec
UPDATE tasks
SET summary = ?,
  assigned_to_user_id = ?,
  finished_at = ?
WHERE id = ?;
-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?;