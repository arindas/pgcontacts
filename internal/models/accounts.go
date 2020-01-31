package models

import (
	"fmt"
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
	Password string `json:"password"`
	Token    string `json:"token" sql:"-"`
}

func init() {
	app.GetDB().AutoMigrate(&Account{})
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

func (account *Account) setToken() {
	tokenWrapper := &app.Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenWrapper)
	if tokenString, err := token.SignedString([]byte(os.Getenv("token_password"))); err != nil {
		panic(err)
	} else {
		account.Token = tokenString
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

	account.setToken()

	response := utils.Message("Account has been created.", true)
	account.Password = ""
	response["account"] = account
	return response
}

func Login(email, password string) map[string]interface{} {

	account := &Account{}
	query := app.GetDB().Table("accounts").Where("email = ?", email).First(account)
	if errors := query.GetErrors(); len(errors) > 0 {
		if query.RecordNotFound() {
			return utils.Message(
				fmt.Sprintf("Account with email: %s not found.", email), false)
		}

		return utils.Message((gorm.Errors)(errors).Error(), false)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return utils.Message("Invalid login credentials, please try again.", false)
		}

		panic(err)
	}

	account.setToken()

	response := utils.Message("Logged in.", true)
	account.Password = ""
	response["account"] = account
	return response
}

func GetUser(id uint) *Account {
	account := &Account{}
	query := app.GetDB().Table("accounts").Where("id = ?", id).First(account)

	if errors := query.GetErrors(); len(errors) > 0 {
		return nil
	}

	account.Password = ""
	return account
}
