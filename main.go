package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rinha-backend/config"
	"time"

	_ "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const port = "8080"

func main() {
	log.Println("Starting authentication service")
	pool, err := connectToDB()
	if err != nil {
		log.Panic("Can't connect to Postgres")
	}

	defer pool.Close()

	app := config.NewConfig(pool)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.Routers,
	}
	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 10; i++ {
		log.Printf("Connecting to Postgres, attempt %d\n", i+1)

		pool, err = pgxpool.ConnectConfig(context.Background(), config.ConfigPGX())
		if err != nil {
			log.Printf("Failed to connect to Postgres: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to Postgres")
		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to Postgres after 10 attempts")
}
