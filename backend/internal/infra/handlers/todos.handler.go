package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mori5321/mangahash/backend/internal/infra/errs"
	"github.com/mori5321/mangahash/backend/internal/usecase"
)

func TodosHandler(stores usecase.Stores) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todos, err := listHandler(r, stores)
			handleResponse(w, todos, http.StatusOK, err)
		case http.MethodPost:
			todo, err := addHandler(r, stores)
			handleResponse(w, todo, http.StatusCreated, err)
		default:
			handleError(w, errs.MethodNotAllowedError)
		}
	}
}

func listHandler(r *http.Request, stores usecase.Stores) ([]usecase.TodoDTO, error) {
	u := usecase.NewTodoUsecase(stores)
	todos, err := u.ListTodos(nil, nil)

	if err != nil {
		return nil, err
	}

	if len(todos) == 0 {
		return []usecase.TodoDTO{}, nil
	}

	return todos, nil
}

func addHandler(r *http.Request, stores usecase.Stores) (*usecase.TodoDTO, error) {
	var input usecase.CreateTodoInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, fmt.Errorf("Invalid Request Body %s", input, errs.InvalidRequestError)
	}

	u := usecase.NewTodoUsecase(stores)
	todo, err := u.AddTodo(input)

	if err != nil {
		return nil, err
	}

	return todo, nil
}

func TodoHandler(stores usecase.Stores) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getParams(r, "/todos")

		switch r.Method {
		case http.MethodGet:
			todo, err := fetchHandler(r, stores)
			handleResponse(w, todo, http.StatusOK, err)
		default:
			handleError(w, fmt.Errorf("Method Not Allowed: %s %w", r.Method, errs.MethodNotAllowedError))
		}
	}
}

func fetchHandler(r *http.Request, stores usecase.Stores) (*usecase.TodoDTO, error) {
	params := getParams(r, "/todos")
	id := params[0]

	u := usecase.NewTodoUsecase(stores)
	todo, err := u.FetchTodoByID(id)

	if err != nil {
		return nil, err
	}

	return todo, nil
}
