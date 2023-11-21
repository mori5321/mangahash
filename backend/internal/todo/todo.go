package todo

type Todo struct {
	ID    string
	Title string
}

func NewTodo(id string, title string) Todo {
	return Todo{
		ID:    id,
		Title: title,
	}
}

type ListPagination struct {
	Limit  int32
	Offset int32
}

type TodoRepository interface {
	List(pagination *ListPagination) ([]Todo, error)
	Fetch(id string) (*Todo, error)
	Store(todo Todo) error
	Delete(id string) error
}

type IDGenerator interface {
	Gen() string
	Valid(s string) bool
}
