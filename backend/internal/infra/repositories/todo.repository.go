package repositories

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mori5321/mangahash/backend/internal/core/todo"
	"github.com/mori5321/mangahash/backend/internal/infra/errs"
	"github.com/mori5321/mangahash/backend/queries"
)

type TodoRepositoryPostgres struct {
	dbConn *pgx.Conn
}

func NewTodoRepositoryPostgres(dbConn *pgx.Conn) *TodoRepositoryPostgres {
	return &TodoRepositoryPostgres{dbConn: dbConn}
}

func (repo *TodoRepositoryPostgres) List(pagination *todo.ListPagination) ([]todo.Todo, error) {
	ctx := context.Background() // ?: context の扱いよくわからない

	query := queries.New(repo.dbConn)

	defaultLimit := int32(30)
	defaultOffset := int32(0)

	if pagination == nil {
		pagination = &todo.ListPagination{
			Limit:  defaultLimit,
			Offset: defaultOffset,
		}
	}

	models, err := query.ListTodos(ctx, queries.ListTodosParams{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
	})

	if err != nil {
		return nil, err
	}

	todos := make([]todo.Todo, len(models))
	for i, model := range models {
		todo := repo.toEntity(model)
		todos[i] = todo
	}

	return todos, nil
}

func (repo *TodoRepositoryPostgres) Store(todo todo.Todo) error {
	return nil
}

func (repo *TodoRepositoryPostgres) Fetch(id string) (*todo.Todo, error) {
	query := queries.New(repo.dbConn)

	uid := uuid.Must(uuid.FromString(id))

	model, err := query.FetchTodo(context.TODO(), uid)
	if err != nil {
		return nil, errs.NotFoundError
	}

	todo := repo.toEntity(model)

	return &todo, nil
}

func (repo *TodoRepositoryPostgres) Delete(id string) error {
	return nil
}

func (repo *TodoRepositoryPostgres) toEntity(todoModel queries.Todo) todo.Todo {
	return todo.NewTodo(todoModel.ID.String(), todoModel.Title)
}
