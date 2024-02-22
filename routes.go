package main

import (
	"net/http"
	"rinha-backend/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *Config) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middlewares.ValidateID)
	// mux.Post("/clientes/:id/transacoes", CreateTransactionHandler)
	mux.Get("/clientes/{id}/extrato", app.GetTransactionsHandler)
	mux.Post("/cadastrar", app.CreateClientHandler)
	mux.Post("/clientes/{id}/transacoes", app.CreateNewTransactionHandler)

	return mux
}
