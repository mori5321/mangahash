package todo

import (
	"errors"

	"github.com/mori5321/mangahash/backend/internal/common"
)

type Stores struct {
	IDGenerator IDGenerator
	TodoRepo    TodoRepository
}

func newStores(idGenerator IDGenerator, todoRepo TodoRepository) Stores {
	return Stores{
		IDGenerator: idGenerator,
		TodoRepo:    todoRepo,
	}
}

type TodoUsecase struct {
	stores Stores
}

func NewTodoUsecase(stores Stores) TodoUsecase {
	return TodoUsecase{stores: stores}
}

// このlayerにjson serializerを書くのはどうかと思うが、
// あまりにlayer間の変換層が多くなるのも非生産的なので、一旦この方針で進む。
type TodoDTO struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type CreateTodoInput struct {
	Title string `json:"title"`
}

func (i CreateTodoInput) Validate() error {
	var errs []error

	if i.Title == "" {
		errs = append(errs, errors.New("title is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (u TodoUsecase) toDTO(todo *Todo) TodoDTO {
	return TodoDTO{
		ID:    todo.ID,
		Title: todo.Title,
	}
}

func (u TodoUsecase) ListTodos(limit *int, offset *int) ([]TodoDTO, error) {
	pagination := &ListPagination{
		Limit:  30,
		Offset: 0,
	}
	todos, err := u.stores.TodoRepo.List(pagination)
	if err != nil {
		return nil, err
	}

	todoDTOs := make([]TodoDTO, len(todos))
	for i, todo := range todos {
		todoDTO := u.toDTO(&todo)
		todoDTOs[i] = todoDTO
	}

	return todoDTOs, nil
}

func (u TodoUsecase) FetchTodoByID(id string) (*TodoDTO, error) {
	if !u.stores.IDGenerator.Valid(id) {
		return nil, common.InvalidIDError
	}

	todo, err := u.stores.TodoRepo.Fetch(id)
	if err != nil {
		return nil, err
	}

	dto := u.toDTO(todo)

	return &dto, nil
}

func (u TodoUsecase) AddTodo(input CreateTodoInput) (*TodoDTO, error) {
	id := u.stores.IDGenerator.Gen()
	newTodo := NewTodo(
		id,
		input.Title,
	)

	err := u.stores.TodoRepo.Store(newTodo)
	if err != nil {
		return nil, err
	}

	dto := u.toDTO(&newTodo)

	return &dto, nil
}
