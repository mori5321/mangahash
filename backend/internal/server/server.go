package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/mori5321/mangahash/backend/internal/healthz"
	"github.com/mori5321/mangahash/backend/internal/todo"
)

func router(dbConn *pgx.Conn) *http.ServeMux {
	// 参考
	// https://ema-hiro.hatenablog.com/entry/2018/10/22/015427
	// https://journal.lampetty.net/entry/understanding-http-handler-in-go

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz.HealthzHandler)
	mux.HandleFunc("/todos", todo.TodosHandler(dbConn))
	mux.HandleFunc("/todos/", todo.TodoHandler(dbConn))

	return mux
}

type AppConfig struct {
	AppPort    uint16
	DBUser     string
	DBName     string
	DBHost     string
	DBPassword string
}

func App(ctx context.Context, conf AppConfig, readyChForTest chan<- bool) {
	dbConn, err := pgx.Connect(ctx, fmt.Sprintf("user=%s dbname=%s host=%s password=%s\n", conf.DBUser, conf.DBName, conf.DBHost, conf.DBPassword))
	defer dbConn.Close(ctx)

	router := router(dbConn)

	if readyChForTest != nil {
		readyChForTest <- true
	}

	fmt.Printf("Server is now running on port %d\n", conf.AppPort)
	port := fmt.Sprintf(":%d", conf.AppPort)
	err = http.ListenAndServe(port, router)

	if err != nil {
		log.Fatal(err)
	}
}

func Run() {
	fmt.Printf("We're launching a server now ... \n")
	ctx := context.Background()
	conf := AppConfig{
		AppPort:    9090,
		DBUser:     "postgres",
		DBName:     "docker",
		DBHost:     "db",
		DBPassword: "password",
	}

	App(ctx, conf, nil)
}
