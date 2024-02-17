package main

import (
	"encoding/json"
	"net/http"
	"rinha-backend/data"
	"rinha-backend/helpers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *Config) CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var client data.Client

	if erro := json.NewDecoder(r.Body).Decode(&client); erro != nil {
		helpers.ErrorJSON(w, erro, http.StatusBadRequest)
		return
	}

	newId, erro := data.Models.CreateClientModel(data.Models{}, client)
	if erro != nil {
		helpers.ErrorJSON(w, erro, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, newId)
}

func (app *Config) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	// como sabemos, para capturar o parametro de busca de uma rota, precisamos primeiro capturar ele de forma dinâmica. parra isso, usamos o query string, que em go é feito da maneira abaixo.
	clientIdStr := chi.URLParam(r, "id")

	// Converter o valor para um inteiro
	clientId, _ := strconv.Atoi(clientIdStr)

	transactions, erro := data.Models.GetTransactionsModel(data.Models{}, clientId)
	// data.Models.GetExtractHandler(*data.Models.Statement{}, clientId)
	if erro != nil {
		helpers.ErrorJSON(w, erro, http.StatusUnauthorized)
	}

	helpers.WriteJSON(w, http.StatusOK, transactions)
}

func (app *Config) CreateNewTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var newTransaction data.Transactions

	if erro := json.NewDecoder(r.Body).Decode(&newTransaction); erro != nil {
		helpers.ErrorJSON(w, erro, http.StatusBadRequest)
		return
	}

	clientIdStr := chi.URLParam(r, "id")

	// Converter o valor para um inteiro
	clientId, _ := strconv.Atoi(clientIdStr)

	transactionResult, erro := data.Models.CreateTransactionModel(data.Models{}, newTransaction, clientId)
	if erro != nil {
		helpers.ErrorJSON(w, erro, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, transactionResult)
}
