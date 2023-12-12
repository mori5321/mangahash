package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/mori5321/mangahash/crawler/internal/healthz"
)

func router() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz.HealthzHandler)

	mux.HandleFunc("/crawler/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Yet to be implemented")
	})

	return mux
}

type AppConfig struct {
	AppPort uint16
}

func App(ctx context.Context, conf AppConfig, readyChForTest chan<- bool) {
	router := router()

	if readyChForTest != nil {
		readyChForTest <- true
	}

	fmt.Printf("Server is now running on port %d\n", conf.AppPort)
	port := fmt.Sprintf(":%d", conf.AppPort)
	err := http.ListenAndServe(port, router)

	if err != nil {
		log.Fatal(err)
	}
}

func Run() {
	fmt.Printf("We're launching a server now ... \n")
	ctx := context.Background()

	conf := AppConfig{
		AppPort: 5050, // TODO: get from env
	}

	App(ctx, conf, nil)
}
