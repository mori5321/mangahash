package todo

import "github.com/mori5321/mangahash/backend/internal/common"

var todos = []Todo{}

type TodoRepositoryMem struct {
}

func NewTodoRepositoryMem() *TodoRepositoryMem {
	return &TodoRepositoryMem{}
}

func (r *TodoRepositoryMem) List(pagination *ListPagination) ([]Todo, error) {
	return todos, nil
}

func (r *TodoRepositoryMem) Fetch(id string) (*Todo, error) {
	for _, todo := range todos {
		if todo.ID == id {
			return &todo, nil
		}
	}
	return nil, common.NotFoundError
}

func (r *TodoRepositoryMem) Store(todo Todo) error {
	for i, t := range todos {
		if t.ID == todo.ID {
			todos[i] = todo
			return nil
		}
	}

	todos = append(todos, todo)
	return nil
}

func (r *TodoRepositoryMem) Delete(id string) error {
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return common.NotFoundError
}
