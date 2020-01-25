package app

import (
	"encoding/json"
	jwt "github.com/dgrijalva/jwt-go"
	"net/http"
	"sort"
	"strings"
)

var authlessPaths = []string{"/api/user/new", "/api/user/login"}

func reportError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Add("Content-Type", "application/json")
	var data map[string]interface{}
	data["status"] = false
	data["messaage"] = err
	json.NewEncoder(w).Encode(data)
}

type Token struct {
	UserId uint
	jwt.StandardClaims
}

func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sort.SearchStrings(authlessPaths, r.URL.Path) > -1 {
			next.ServeHTTP(w, r)
		} else {
			tokenHeader := r.Header.Get("Authorization")

			if len(tokenHeader) == 0 {
				reportError(w, "Missing auth token")
				return
			}

			tokenContents := strings.Split(tokenHeader, " ")
			if len(tokenContents) != 2 {
				reportError(w, "Malformed auth token")
				return
			}

			tokenPart := tokenContents[1]
			tk := &Token{}

		}
	})
}
