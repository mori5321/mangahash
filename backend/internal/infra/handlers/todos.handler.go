package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mori5321/mangahash/backend/queries"
)

type TodoResponse struct {
	Id    string `json:"id"`
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

type TodosResponse = []TodoResponse

func TodosHandler(dbConn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			listHandler(w, r, dbConn)
		case http.MethodPost:
			addHandler(w, r, dbConn)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func listHandler(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
	todos, err := queries.New(dbConn).ListTodos(context.Background(), queries.ListTodosParams{
		Limit:  10,
		Offset: 0,
	})

	if err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode:     InteralServerError,
			ErrorMessages: []string{fmt.Sprintf("Failed to list todos: %s", err)},
		})
		return
	}

	if len(todos) == 0 {
		respondWithJson(w, TodosResponse{}, http.StatusOK)
		return
	}

	respondWithJson(w, todos, http.StatusOK)
}

func addHandler(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
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

	newTodo, err := queries.New(dbConn).CreateTodo(
		context.Background(),
		req.Title,
	)

	if err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode:     InteralServerError,
			ErrorMessages: []string{fmt.Sprintf("Failed to create todo: %s", err)},
		})
		return
	}

	respondWithJson(w, newTodo, http.StatusOK)
}

func TodoHandler(dbConn *pgx.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getParams(r, "/todos")

		switch r.Method {
		case http.MethodGet:
			getHandler(w, r, dbConn)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func getHandler(w http.ResponseWriter, r *http.Request, dbConn *pgx.Conn) {
	params := getParams(r, "/todos")
	var id pgtype.UUID
	if err := id.Scan(params[0]); err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode:     InvalidRequest,
			ErrorMessages: []string{fmt.Sprintf("Failed to parse UUID: %s", err)},
		})
		return
	}

	todo, err := queries.New(dbConn).GetTodo(
		context.Background(),
		id,
	)

	if err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode:     NotFound,
			ErrorMessages: []string{fmt.Sprintf("Failed to get todo: %s", err)},
		})
		return
	}

	idStr, err := todo.ID.Value()
	if err != nil {
		respondWithError(w, ErrorResponse{
			ErrorCode:     InteralServerError,
			ErrorMessages: []string{fmt.Sprintf("Failed to get todo: %s", err)},
		})
	}

	todoResponse := TodoResponse{
		Id:    idStr.(string),
		Title: todo.Title,
	}

	respondWithJson(w, todoResponse, http.StatusOK)
}
