package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/arindas/pgcontacts/internal/models"
	"github.com/arindas/pgcontacts/internal/utils"
)

func decodeAccountJson(w http.ResponseWriter, r *http.Request) *models.Account {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		utils.Respond(w, utils.Message("Invalid request.", false))
		return nil
	}

	return account
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if account := decodeAccountJson(w, r); account != nil {
		response := account.Create()
		utils.Respond(w, response)
	}
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	if account := decodeAccountJson(w, r); account != nil {
		response := models.Login(account.Email, account.Password)
		utils.Respond(w, response)
	}
}
