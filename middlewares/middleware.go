package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"rinha-backend/helpers"

	"github.com/go-chi/chi/v5"
)

func ValidateID(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, err := extractID(r)

			if err != nil {
				helpers.ErrorJSON(w, err, http.StatusUnprocessableEntity)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}

func extractID(r *http.Request) (string, error) {
	id := chi.URLParam(r, "id")

	fmt.Println(r)
	if id == "" {
		return "", errors.New("cliente não encontrado para o ID fornecido")
	}
	return id, nil
}
