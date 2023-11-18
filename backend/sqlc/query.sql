-- name: ListTodos :many
SELECT id, title FROM todos LIMIT $1 OFFSET $2;

-- name: FetchTodo :one
SELECT id, title FROM todos WHERE id = $1;

-- name: CreateTodo :exec
INSERT INTO todos (title) VALUES ($1);

-- name: UpdateTodo :exec
UPDATE todos SET title = $1 WHERE id = $2;

-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1;

