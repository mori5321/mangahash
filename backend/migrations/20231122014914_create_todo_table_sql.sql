-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS app.todos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS app.todos;
-- +goose StatementEnd

