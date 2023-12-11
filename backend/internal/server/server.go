package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mori5321/mangahash/backend/internal/healthz"
	"github.com/mori5321/mangahash/backend/internal/todo"
	"github.com/mori5321/mangahash/backend/queries"
)

func router(dbPool queries.DBTX) *http.ServeMux {
	// 参考
	// https://ema-hiro.hatenablog.com/entry/2018/10/22/015427
	// https://journal.lampetty.net/entry/understanding-http-handler-in-go

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthz.HealthzHandler)
	mux.HandleFunc("/todos", todo.TodosHandler(dbPool))
	mux.HandleFunc("/todos/", todo.TodoHandler(dbPool))

	return mux
}

type AppConfig struct {
	AppPort    uint16
	DBUser     string
	DBName     string
	DBHost     string
	DBPort     uint16
	DBPassword string
}

func connectDBPool(ctx context.Context, conf AppConfig) (*pgxpool.Pool, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s host=%s port=%d password=%s\n", conf.DBUser, conf.DBName, conf.DBHost, conf.DBPort, conf.DBPassword)
	fmt.Printf("connStr: %s\n", connStr)
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

func App(ctx context.Context, conf AppConfig, readyChForTest chan<- bool) {
	pool, err := connectDBPool(ctx, conf)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect DB: %w", err))
	}
	defer pool.Close()

	router := router(pool)

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

	dbhost := os.Getenv("DATABASE_HOST")
	dbpass := os.Getenv("DATABASE_PASSWORD")
	dbname := os.Getenv("DATABASE_NAME")
	dbuser := os.Getenv("DATABASE_USER")

	tmpdbport, err := strconv.ParseUint(os.Getenv("DATABASE_PORT"), 10, 16)
	if err != nil {
		log.Fatal(err)
	}
	dbport := uint16(tmpdbport)

	conf := AppConfig{
		AppPort:    9090, // TODO: get from env
		DBUser:     dbuser,
		DBName:     dbname,
		DBPort:     dbport,
		DBHost:     dbhost,
		DBPassword: dbpass,
	}

	App(ctx, conf, nil)
}
