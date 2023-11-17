package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mori5321/mangahash/backend/internal/server/handlers"
)

func Router() *http.ServeMux {
	// 参考
	// https://ema-hiro.hatenablog.com/entry/2018/10/22/015427
	// https://journal.lampetty.net/entry/understanding-http-handler-in-go

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", handlers.HealthzHandler)
	mux.HandleFunc("/todos", handlers.TodosHandler)
	mux.HandleFunc("/todos/", handlers.TodoHandler)

	return mux
}

func Run() {
	router := Router()

	fmt.Println("Server is now running on port 9090")
	err := http.ListenAndServe(":9090", router)

	if err != nil {
		log.Fatal(err)
	}

}
