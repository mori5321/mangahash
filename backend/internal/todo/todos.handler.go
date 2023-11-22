package todo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/mori5321/mangahash/backend/internal/common"
)

func TodosHandler(dbConn *pgx.Conn) http.HandlerFunc {
	idGenerator := common.NewUUIDGenerator()
	todoRepo := NewTodoRepositoryPostgres(dbConn)
	stores := newStores(idGenerator, todoRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todos, err := listHandler(r, stores)
			common.HandleResponse(w, todos, http.StatusOK, err)
		case http.MethodPost:
			todo, err := addHandler(r, stores)
			fmt.Println(todo, err)
			common.HandleResponse(w, todo, http.StatusCreated, err)
		default:
			common.HandleError(w, common.MethodNotAllowedError)
		}
	}
}

func listHandler(r *http.Request, stores Stores) ([]TodoDTO, error) {
	u := NewTodoUsecase(stores)
	todos, err := u.ListTodos(nil, nil)

	if err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		return []TodoDTO{}, nil
	}

	return todos, nil
}

func addHandler(r *http.Request, stores Stores) (*TodoDTO, error) {
	var input CreateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fmt.Errorf("Invalid Request Body %s %w", input, common.InvalidRequestError)
	}

	u := NewTodoUsecase(stores)
	todo, err := u.AddTodo(input)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func TodoHandler(dbConn *pgx.Conn) http.HandlerFunc {
	idGenerator := common.NewUUIDGenerator()
	todoRepo := NewTodoRepositoryPostgres(dbConn)
	stores := newStores(idGenerator, todoRepo)

	return func(w http.ResponseWriter, r *http.Request) {
		common.GetParams(r, "/todos")

		switch r.Method {
		case http.MethodGet:
			todo, err := fetchHandler(r, stores)
			common.HandleResponse(w, todo, http.StatusOK, err)
		default:
			common.HandleError(w, fmt.Errorf("Method Not Allowed: %s %w", r.Method, common.MethodNotAllowedError))
		}
	}
}

func fetchHandler(r *http.Request, stores Stores) (*TodoDTO, error) {
	params := common.GetParams(r, "/todos")
	id := params[0]

	u := NewTodoUsecase(stores)
	todo, err := u.FetchTodoByID(id)

	if err != nil {
		return nil, err
	}

	return todo, nil
}
