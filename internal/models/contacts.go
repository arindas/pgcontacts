package models

import (
	"github.com/jinzhu/gorm"

	"github.com/arindas/pgcontacts/internal/app"
	"github.com/arindas/pgcontacts/internal/utils"
)

type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserId uint   `json:"user_id" gorm:"user_id"`
}

func init() {
	app.GetDB().AutoMigrate(&Contact{})
}

func (contact *Contact) Validate() (map[string]interface{}, bool) {

	if len(contact.Name) == 0 || len(contact.Phone) == 0 {
		return utils.Message("Missing Name or Phone number in payload.", false), false
	}

	if contact.UserId <= 0 || GetUser(contact.UserId) == nil {
		return utils.Message("Invalid or unregistered user ID.", false), false
	}

	return utils.Message("Valid contact", true), true
}

func (contact *Contact) Create() map[string]interface{} {
	if response, ok := contact.Validate(); !ok {
		return response
	}

	app.GetDB().Create(contact)

	response := utils.Message("Contact created.", true)
	response["contact"] = contact
	return response
}

func GetContactById(id uint) *Contact {
	contact := &Contact{}
	query := app.GetDB().Table("contacts").Where(
		"id = ?", id).First(contact)

	if errors := query.GetErrors(); len(errors) > 0 {
		return nil
	}

	return contact
}

func GetContactsForUser(userId uint) []*Contact {
	contacts := make([]*Contact, 0)
	query := app.GetDB().Table("contacts").Where(
		"user_id = ?", userId).Find(&contacts)

	if errors := query.GetErrors(); len(errors) > 0 {
		return nil
	}

	return contacts
}
