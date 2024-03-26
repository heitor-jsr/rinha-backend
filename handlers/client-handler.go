package handlers

import (
	"encoding/json"
	"net/http"
	"rinha-backend/data"
	"rinha-backend/helpers"
)

func CreateClientHandler(w http.ResponseWriter, r *http.Request) {

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
