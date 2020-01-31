package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/arindas/pgcontacts/internal/models"
	"github.com/arindas/pgcontacts/internal/utils"
)

func CreateContact(w http.ResponseWriter, r *http.Request) {
	contact := &models.Contact{}
	idInContext := r.Context().Value("user")
	if idInContext == nil {
		utils.Respond(w, utils.Message("User Id not found in request context.", false))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(contact); err != nil {
		utils.Respond(w, utils.Message("Error decoding contact json.", false))
		return
	}

	contact.UserId = idInContext.(uint)
	response := contact.Create()

	response["contact"] = contact
	utils.Respond(w, response)
}

func GetContactsForUser(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)

	contacts := models.GetContactsForUser(userId)
	var response map[string]interface{}

	if contacts != nil {
		response = utils.Message("Contacts retrieved.", true)
		response["contacts"] = contacts
	} else {
		response = utils.Message("No contacts retrieved", false)
	}
	utils.Respond(w, response)
}
