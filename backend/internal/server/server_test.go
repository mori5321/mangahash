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
)

var conf AppConfig = AppConfig{
	AppPort:    9091,
	DBUser:     "test",
	DBName:     "test",
	DBHost:     "testdb",
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
	// Send a request to the server
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	// Check the status code is what we expect.
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestTodosHandler(t *testing.T) {
	url := appHost + "/todos"

	td := todo.CreateTodoInput{
		Title: "test",
	}
	j, err := json.Marshal(td)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(j)

	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		panic(err)
	}

	res, err = http.Get(url)
	if err != nil {
		panic(err)
	}

	status := res.StatusCode
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := res.Body
	defer body.Close()

	// Assert response body is JSON empty array
	var todos []todo.TodoDTO
	err = json.NewDecoder(body).Decode(&todos)
	if err != nil {
		panic(err)
	}

	// TODO: Need test db and reset system for response body test
	if len(todos) != 1 {
		t.Errorf("GET /todos is supposed to return array with length = %d: got %d", 1, len(todos))
	}
}
