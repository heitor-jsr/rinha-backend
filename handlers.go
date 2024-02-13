package main

import (
	"errors"
	"fmt"
	"net/http"
	"rinha-backend/data"
	"rinha-backend/helpers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Config) GetTransactions(w http.ResponseWriter, r *http.Request) {
	// como sabemos, para capturar o parametro de busca de uma rota, precisamos primeiro capturar ele de forma dinâmica. parra isso, usamos o query string, que em go é feito da maneira abaixo.
	clientIdStr := chi.URLParam(r, "id")
	fmt.Println(clientIdStr)

	// Converter o valor para um inteiro
	clientId, _ := strconv.Atoi(clientIdStr)

	transactions, erro := data.Models.GetExtractHandler(data.Models{}, clientId)
	// data.Models.GetExtractHandler(*data.Models.Statement{}, clientId)
	if erro != nil {
		helpers.ErrorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
	}

	helpers.WriteJSON(w, http.StatusOK, transactions)
}
