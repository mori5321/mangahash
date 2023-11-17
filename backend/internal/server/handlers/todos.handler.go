package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type TodoResponse struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

type CreateTodoRequest struct {
	Title string `json:"title"`
}

func (r CreateTodoRequest) Validate() error {
	var errs []error

	if r.Title == "" {
		errs = append(errs, errors.New("title is required"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

var todos = []TodoResponse{
	{Id: 1, Title: "Buy milk"},
	{Id: 2, Title: "Buy eggs"},
	{Id: 3, Title: "Buy bread"},
}

type TodosResponse = []TodoResponse

func TodosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listHandler(w, r)
	case http.MethodPost:
		addHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJson(w, todos, http.StatusOK)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode: InvalidRequest,
			ErrorMessages: []string{
				fmt.Sprintf("Failed to decode request body: %s", err),
			},
		})
		return
	}

	if err := req.Validate(); err != nil {
		var msgs []string
		if errs, ok := err.(interface{ Unwrap() []error }); ok {
			for _, err := range errs.Unwrap() {
				msgs = append(msgs, err.Error())
			}
		}

		if len(msgs) == 0 {
			msgs = append(msgs, "Failed to validate request")
		}

		respondWithError(w, ErrorResponse{
			ErrorCode:     InvalidRequest,
			ErrorMessages: msgs,
		})
		return
	}

	newTodo := TodoResponse{
		Id:    len(todos) + 1,
		Title: req.Title,
	}

	todos = append(todos, newTodo)

	respondWithJson(w, newTodo, http.StatusOK)
}

func TodoHandler(w http.ResponseWriter, r *http.Request) {
	getParams(r, "/todos")

	switch r.Method {
	case http.MethodGet:
		getHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	params := getParams(r, "/todos")
	id, err := strconv.Atoi(params[0])

	if err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode:     InvalidRequest,
			ErrorMessages: []string{"id must be an integer"},
		})
		return
	}

	var t TodoResponse
	for i, todo := range todos {
		if todo.Id == id {
			t = todo
			break
		}

		if i == len(todos)-1 {
			respondWithError(w, ErrorResponse{
				ErrorCode:     NotFound,
				ErrorMessages: []string{fmt.Sprintf("Todo with id %d not found", id)},
			})
			return
		}
	}

	respondWithJson(w, t, http.StatusOK)
}
