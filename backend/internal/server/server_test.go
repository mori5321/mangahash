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

	// POST /todos ----------
	td := todo.CreateTodoInput{
		Title: "test",
	}
	j, err := json.Marshal(td)
	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(j)
	_, err = http.Post(url, "application/json", reader)
	if err != nil {
		panic(err)
	}

	// GET /todos ----------
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	status := res.StatusCode
	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	body := res.Body
	defer body.Close()

	var todos []todo.TodoDTO
	err = json.NewDecoder(body).Decode(&todos)
	if err != nil {
		panic(err)
	}

	if len(todos) != 1 {
		t.Errorf("GET /todos is supposed to return array with length = %d: got %d", 1, len(todos))
	}

	// GET /todos/:id --------
	first := todos[0]
	res, err = http.Get(url + "/" + first.ID)
	if err != nil {
		panic(err)
	}

	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var result todo.TodoDTO
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	if result.Title != td.Title {
		t.Errorf("GET /todos/:id is supposed to return title = %s: got %s", td.Title, result.Title)
	}

	// PUT /todos/:id --------
	updated := todo.UpdateTodoInput{
		Title: "updated",
	}

	j, err = json.Marshal(updated)
	if err != nil {
		panic(err)
	}

	reader = bytes.NewReader(j)
	req, err := http.NewRequest(http.MethodPut, url+"/"+first.ID, reader)
	if err != nil {
		panic(err)
	}

	http.DefaultClient.Do(req)

	res, err = http.Get(url + "/" + first.ID)
	if err != nil {
		panic(err)
	}

	if status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		panic(err)
	}

	if result.Title != updated.Title {
		t.Errorf("GET /todos/:id is supposed to return title = %s: got %s", td.Title, result.Title)
	}

	// DELETE /todos/:id -----
	req, err = http.NewRequest(http.MethodDelete, url+"/"+first.ID, nil)
	if err != nil {
		panic(err)
	}

	res, err = http.DefaultClient.Do(req)

	if res.StatusCode != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	res, err = http.Get(url + "/" + first.ID)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	res, err = http.Get(url)
	if err != nil {
		panic(err)
	}

	body = res.Body
	defer body.Close()

	err = json.NewDecoder(body).Decode(&todos)
	if err != nil {
		panic(err)
	}

	if len(todos) != 0 {
		t.Errorf("GET /todos is supposed to return array with length = %d: got %d", 0, len(todos))
	}
}
