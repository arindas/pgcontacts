package models

import (
	"github.com/jinzhu/gorm"
)

// Account represents user accounts
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

func (account *Account) Validate() (map[string]interface{}, bool) {
	var message map[string]interface{}
	return message, false
}
