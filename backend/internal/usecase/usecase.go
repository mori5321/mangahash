package usecase

import (
	"github.com/mori5321/mangahash/backend/internal/core/todo"
)

type Stores struct {
	TodoRepo todo.TodoRepository
	UUID     UUIDGenerator
}

func NewStores(uuid UUIDGenerator, todoRepo todo.TodoRepository) Stores {
	return Stores{
		UUID:     uuid,
		TodoRepo: todoRepo,
	}
}

type UUIDGenerator interface {
	Gen() string
	Valid(s string) bool
}
