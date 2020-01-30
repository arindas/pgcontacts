package models

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"os"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/arindas/pgcontacts/internal/app"
	"github.com/arindas/pgcontacts/internal/utils"
)

// Account represents user accounts
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"-"`
	Token    string `json:"token" sql:"-"`
}

func validEmail(email string) bool {
	emailRegex := `^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		panic(err)
	}

	return matched
}

// Validate returns a json encodable message and a bool
// to state whether the given account is valid or not
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !validEmail(account.Email) {
		return utils.Message("Invalid email", false), false
	}

	if len(account.Password) < 6 {
		return utils.Message("Password length < 6", false), false
	}

	dupEmailAccount := &Account{}
	dupEmailQuery := app.GetDB().Table("accounts").Where("email = ?",
		account.Email).First(dupEmailAccount)

	if dupEmailQuery.RecordNotFound() {
		return utils.Message("Valid account", true), true
	} else if errors := dupEmailQuery.GetErrors(); len(errors) > 0 {
		return utils.Message((gorm.Errors)(errors).Error(), false), false
	} else {
		return utils.Message("Account with duplicate email found.", false), false
	}
}

// Creates an account and returns a json encodable response
// containing the details of the account creation
func (account *Account) Create() map[string]interface{} {
	if response, ok := account.Validate(); !ok {
		return response
	}

	if hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(account.Password), bcrypt.DefaultCost); err != nil {
		panic(err)
	} else {
		account.Password = string(hashedPassword)
	}

	app.GetDB().Create(account)

	if account.ID <= 0 {
		return utils.Message("Failed to create account", false)
	}

	tokenWrapper := &app.Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenWrapper)
	if tokenString, err := token.SignedString([]byte(os.Getenv("token_password"))); err != nil {
		panic(err)
	} else {
		account.Token = tokenString
	}

	response := utils.Message("Account has been created.", true)
	response["account"] = account
	return response
}
