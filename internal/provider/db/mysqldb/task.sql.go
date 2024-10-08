// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: task.sql

package mysqldb

import (
	"context"
	"database/sql"
)

const createTask = `-- name: CreateTask :exec
INSERT INTO tasks (summary, created_by_user_id, assigned_to_user_id)
VALUES (?, ?, ?)
`

type CreateTaskParams struct {
	Summary          string
	CreatedByUserID  string
	AssignedToUserID sql.NullString
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) error {
	_, err := q.db.ExecContext(ctx, createTask, arg.Summary, arg.CreatedByUserID, arg.AssignedToUserID)
	return err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?
`

func (q *Queries) DeleteTask(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteTask, id)
	return err
}

const getTaskByID = `-- name: GetTaskByID :one
SELECT id, summary, assigned_to_user_id, created_by_user_id, finished_at, created_at, updated_at
FROM tasks
WHERE id = ?
LIMIT 1
`

func (q *Queries) GetTaskByID(ctx context.Context, id string) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTaskByID, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Summary,
		&i.AssignedToUserID,
		&i.CreatedByUserID,
		&i.FinishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listTasks = `-- name: ListTasks :many
SELECT id, summary, assigned_to_user_id, created_by_user_id, finished_at, created_at, updated_at
FROM tasks
`

func (q *Queries) ListTasks(ctx context.Context) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Summary,
			&i.AssignedToUserID,
			&i.CreatedByUserID,
			&i.FinishedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTasksWithFilters = `-- name: ListTasksWithFilters :many
SELECT id, summary, assigned_to_user_id, created_by_user_id, finished_at, created_at, updated_at
FROM tasks
WHERE assigned_to_user_id = ?
`

func (q *Queries) ListTasksWithFilters(ctx context.Context, assignedToUserID sql.NullString) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, listTasksWithFilters, assignedToUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Summary,
			&i.AssignedToUserID,
			&i.CreatedByUserID,
			&i.FinishedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :exec
UPDATE tasks
SET summary = ?,
  assigned_to_user_id = ?,
  finished_at = ?
WHERE id = ?
`

type UpdateTaskParams struct {
	Summary          string
	AssignedToUserID sql.NullString
	FinishedAt       sql.NullTime
	ID               string
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) error {
	_, err := q.db.ExecContext(ctx, updateTask,
		arg.Summary,
		arg.AssignedToUserID,
		arg.FinishedAt,
		arg.ID,
	)
	return err
}
