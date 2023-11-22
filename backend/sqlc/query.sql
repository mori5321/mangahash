-- name: ListTodos :many
SELECT id, title FROM app.todos LIMIT $1 OFFSET $2;

-- name: FetchTodo :one
SELECT id, title FROM app.todos WHERE id = $1;

-- name: CreateTodo :exec
INSERT INTO app.todos (title) VALUES ($1);

-- name: UpdateTodo :exec
UPDATE app.todos SET title = $1 WHERE id = $2;

-- name: DeleteTodo :exec
DELETE FROM app.todos WHERE id = $1;

