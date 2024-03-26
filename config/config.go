package config

import (
	"log"
	"net/http"
	"os"
	"rinha-backend/data"
	"rinha-backend/routes"
	"time"

	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ConfigPGX() *pgxpool.Config {
	const defaultMaxConns = int32(5)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

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

	return dbConfig
}

type Config struct {
	DB      *pgxpool.Pool
	Models  data.Models
	Routers http.Handler // Use http.Handler directly if your router implements it
}

func NewConfig(db *pgxpool.Pool) *Config {
	return &Config{
		DB:      db,
		Models:  data.New(db),
		Routers: routes.Routers(), // Assuming routes.Routers() returns an http.Handler
	}
}
