package middlewares

import (
	"errors"
	"net/http"
	"rinha-backend/helpers"

	"github.com/gorilla/mux"
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
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		return "", errors.New("ID não encontrado na URL")
	}
	return id, nil
}
