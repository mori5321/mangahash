package internal

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

func RunApp() {
	ctx := context.Background()
	// From local
	// dbConn, err := pgx.Connect(ctx, "user=postgres dbname=docker sslmode=disable host=localhost port=5430 password=password")
	// For docker-compose
	dbConn, err := pgx.Connect(ctx, "user=postgres dbname=docker sslmode=disable host=db password=password")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(ctx)

	router := router(dbConn)

	fmt.Println("Server is now running on port 9090")
	err = http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatal(err)
	}
}
