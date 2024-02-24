package main

import (
	"encoding/json"
	"errors"
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
	if clientIdStr == "" {
		helpers.ErrorJSON(w, errors.New("forneça um ID para a consulta"), http.StatusNotFound)
		return
	}
	// Converter o valor para um inteiro
	clientId, _ := strconv.Atoi(clientIdStr)

	err := errors.New("cliente não encontrado")
	err2 := errors.New("a transação de débito deixaria o saldo inconsistente")
	transactions, erro := data.Models.GetTransactionsModel(data.Models{}, clientId)
	if erro != nil {
		if errors.Is(erro, err) {
			helpers.ErrorJSON(w, erro, http.StatusNotFound)
			return
		}
		if errors.Is(erro, err2) {
			helpers.ErrorJSON(w, erro, http.StatusUnprocessableEntity)
			return
		}
		helpers.ErrorJSON(w, erro, http.StatusInternalServerError)
		return
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
	if clientIdStr == "" {
		helpers.ErrorJSON(w, errors.New("forneça um ID para a consulta"), http.StatusNotFound)
		return
	}

	// Converter o valor para um inteiro
	clientId, _ := strconv.Atoi(clientIdStr)
	err := errors.New("cliente não encontrado")
	err2 := errors.New("a transação de débito deixaria o saldo inconsistente")

	transactionResult, erro := data.Models.CreateTransactionModel(data.Models{}, newTransaction, clientId)
	if erro != nil {
		if errors.Is(erro, err) {
			helpers.ErrorJSON(w, erro, http.StatusNotFound)
			return
		}
		if errors.Is(erro, err2) {
			helpers.ErrorJSON(w, erro, http.StatusUnprocessableEntity)
			return
		}
		helpers.ErrorJSON(w, erro, http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, transactionResult)
}
