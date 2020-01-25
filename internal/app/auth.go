package app

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
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

// Token jwt token wrapper
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// MiddleWare mux middleware for parsing JWT authenctication token
func MiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (sort.SearchStrings(authlessPaths, r.URL.PATH)) {
			tokenHeader := r.Header.Get("Authorization")

			if len(tokenHeader) == 0 {
				reportError(w, "Missing auth token")
				return
			}

			tokenContents := strings.Split(tokenHeader, " ")
			if len(tokenContents) != 2 {
				reportError(w, "Invalid/Malformed auth token")
				return
			}

			tokenPart := tokenContents[1]
			tk := &Token{}

			token, err := jwt.ParseWithClaims(tokenPart, tk, func(token * jwt.Token)) {
				return []byte(os.Getenv("token_password")), nil
			}

			if err != nil {
				reportError(w, "Malformed auth token")
				return
			}

			if !token.Valid {
				reportError(w, "Token is not valid.")
				return
			}

			fmt.Sprintf("User %s", tk.Username)
			ctx := context.WithValue(r.Context(), "user", tk.UserID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
