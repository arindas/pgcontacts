package app

import (
	"context"
	"net/http"
	"os"
	"sort"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/arindas/pgcontacts/internal/utils"
)

var authlessPaths = []string{"/api/user/new", "/api/user/login"}

func respondWithError(w http.ResponseWriter, err string) {
	w.WriteHeader(http.StatusForbidden)
	utils.Respond(w, utils.Message(err, false))
}

// Token is a jwt.StandardClaims wrapper
type Token struct {
	jwt.StandardClaims
	UserID uint
}

// AuthMiddleWare is a mux.MiddleWare for parsing
// JWT authenctication tokens
func AuthMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sort.SearchStrings(authlessPaths, r.URL.Path) < 0 {
			tokenHeader := r.Header.Get("Authorization")

			if len(tokenHeader) == 0 {
				respondWithError(w, "Missing auth token")
				return
			}

			tokenContents := strings.Split(tokenHeader, " ")
			if len(tokenContents) != 2 {
				respondWithError(w, "Invalid/Malformed auth token")
				return
			}

			tokenPart := tokenContents[1]
			tk := &Token{}

			token, err := jwt.ParseWithClaims(tokenPart, tk,
				func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("token_password")), nil
				})

			if err != nil {
				respondWithError(w, "Malformed auth token")
				return
			}

			if !token.Valid {
				respondWithError(w, "Token is not valid.")
				return
			}

			ctx := context.WithValue(r.Context(), "user", tk.UserID)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}
