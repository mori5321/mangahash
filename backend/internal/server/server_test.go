package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/mori5321/mangahash/backend/internal/todo"
	"github.com/stretchr/testify/require"
)

var conf AppConfig = AppConfig{
	AppPort:    9091,
	DBUser:     "test",
	DBName:     "test",
	DBHost:     "testdb",
	DBPort:     5432,
	DBPassword: "password",
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	readyCh := make(chan bool)
	go func() {
		App(ctx, conf, readyCh)
	}()

	<-readyCh

	code := m.Run()

	// ctx.Done()
	os.Exit(code)
}

var appHost string = fmt.Sprintf("http://localhost:%d", conf.AppPort)

func TestHealthzHandler(t *testing.T) {
	url := appHost + "/healthz"

	res, err := http.Get(url)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestTodosHandler(t *testing.T) {
	url := appHost + "/todos"

	// POST /todos ----------
	td := todo.CreateTodoInput{
		Title: "test",
	}
	j, err := json.Marshal(td)
	require.NoError(t, err)

	reader := bytes.NewReader(j)
	_, err = http.Post(url, "application/json", reader)
	require.NoError(t, err)

	// GET /todos ----------
	res, err := http.Get(url)
	require.NoError(t, err)

	status := res.StatusCode
	require.Equal(t, http.StatusOK, status)

	body := res.Body
	defer body.Close()

	var todos []todo.TodoDTO
	err = json.NewDecoder(body).Decode(&todos)
	require.NoError(t, err)
	require.Equal(t, 1, len(todos))

	// GET /todos/:id --------
	first := todos[0]
	res, err = http.Get(url + "/" + first.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var result todo.TodoDTO
	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, first.ID, result.ID)

	// PUT /todos/:id --------
	updateTodo := todo.UpdateTodoInput{
		Title: "updated",
	}
	j, err = json.Marshal(updateTodo)
	require.NoError(t, err)
	reader = bytes.NewReader(j)
	req, err := http.NewRequest(http.MethodPut, url+"/"+first.ID, reader)
	require.NoError(t, err)

	http.DefaultClient.Do(req)

	res, err = http.Get(url + "/" + first.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	err = json.NewDecoder(res.Body).Decode(&result)
	require.NoError(t, err)
	require.Equal(t, updateTodo.Title, result.Title)

	// DELETE /todos/:id -----
	req, err = http.NewRequest(http.MethodDelete, url+"/"+first.ID, nil)
	require.NoError(t, err)

	res, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get(url + "/" + first.ID)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	res, err = http.Get(url)
	require.NoError(t, err)

	body = res.Body
	defer body.Close()

	err = json.NewDecoder(body).Decode(&todos)
	require.NoError(t, err)
	require.Equal(t, 0, len(todos))
}
