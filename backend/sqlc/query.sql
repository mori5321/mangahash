-- name: ListTodos :many
SELECT id, title FROM todos LIMIT $1 OFFSET $2;

-- name: GetTodo :one
SELECT id, title FROM todos WHERE id = $1;

-- name: CreateTodo :one
INSERT INTO todos (title) VALUES ($1) RETURNING id, title;


