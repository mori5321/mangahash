package core

type Todo struct {
	Id    string
	Title string
}

type TodoRepository interface {
	List() ([]Todo, error)
	Fetch(id string) (Todo, error)
	Store(todo Todo) error
	Delete(id string) error
}
