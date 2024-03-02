package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"rinha-backend/data"
	"time"

	_ "github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func ConfigPGX() *pgxpool.Config {
	const defaultMaxConns = int32(5)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	// Your own Database URL
	DATABASE_URL := os.Getenv("DSN")

	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	dbConfig.BeforeAcquire = func(ctx context.Context, c *pgx.Conn) bool {
		log.Println("Before acquiring the connection pool to the database!!")
		return true
	}

	dbConfig.AfterRelease = func(c *pgx.Conn) bool {
		log.Println("After releasing the connection pool to the database!!")
		return true
	}

	return dbConfig
}

const port = "8080"

type Config struct {
	DB     *pgxpool.Pool
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")
	pool, err := connectToDB()
	if err != nil {
		log.Panic("Can't connect to Postgres")
	}

	defer pool.Close()

	app := Config{
		DB:     pool,
		Models: data.New(pool),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
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

		pool, err = pgxpool.ConnectConfig(context.Background(), ConfigPGX())
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
