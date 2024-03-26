package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"rinha-backend/data"
	"rinha-backend/helpers"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {

	clientIdStr := chi.URLParam(r, "id")
	if clientIdStr == "" {
		helpers.ErrorJSON(w, errors.New("forneça um ID para a consulta"), http.StatusNotFound)
		return
	}
	clientId, _ := strconv.Atoi(clientIdStr)

	transactions, erro := data.Models.GetTransactionsModel(data.Models{}, clientId)
	if erro != nil {
		switch erro.Error() {
		case "cliente não encontrado":
			helpers.ErrorJSON(w, erro, http.StatusNotFound)
		case "a transação de débito deixaria o saldo inconsistente", "tipo de transação inválido", "descrição deve ter entre 1 e 10 caracteres":
			helpers.ErrorJSON(w, erro, http.StatusUnprocessableEntity)
		default:
			helpers.ErrorJSON(w, erro, http.StatusInternalServerError)
		}
		return
	}

	helpers.WriteJSON(w, http.StatusOK, transactions)
}

func CreateNewTransactionHandler(w http.ResponseWriter, r *http.Request) {

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

	clientId, _ := strconv.Atoi(clientIdStr)

	transactionResult, erro := data.Models.CreateTransactionModel(data.Models{}, newTransaction, clientId)
	if erro != nil {
		switch erro.Error() {
		case "cliente não encontrado":
			helpers.ErrorJSON(w, erro, http.StatusNotFound)
		case "a transação de débito deixaria o saldo inconsistente", "tipo de transação inválido", "descrição deve ter entre 1 e 10 caracteres":
			helpers.ErrorJSON(w, erro, http.StatusUnprocessableEntity)
		default:
			helpers.ErrorJSON(w, erro, http.StatusInternalServerError)
		}
		return
	}

	helpers.WriteJSON(w, http.StatusOK, transactionResult)
}
