package infra

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/mori5321/mangahash/backend/internal/infra/handlers"
)

func Router(dbConn *pgx.Conn) *http.ServeMux {
	// 参考
	// https://ema-hiro.hatenablog.com/entry/2018/10/22/015427
	// https://journal.lampetty.net/entry/understanding-http-handler-in-go

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handlers.HealthzHandler)
	mux.HandleFunc("/todos", handlers.TodosHandler(dbConn))
	mux.HandleFunc("/todos/", handlers.TodoHandler(dbConn))

	return mux
}

func Run() {
	ctx := context.Background()
	dbConn, err := pgx.Connect(ctx, "user=postgres dbname=docker sslmode=disable host=db password=password")
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(ctx)

	router := Router(dbConn)

	fmt.Println("Server is now running on port 9090")
	err = http.ListenAndServe(":9090", router)
	if err != nil {
		log.Fatal(err)
	}

}
