package main

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"net/http"

	"github.com/arindas/pgcontacts/internal/app"
	"github.com/arindas/pgcontacts/internal/controllers"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.AuthMiddleWare)

	router.HandleFunc("/api/user/new",
		controllers.CreateAccount).Methods("POST")

	router.HandleFunc("/api/user/login",
		controllers.Authenticate).Methods("POST")

	port := os.Getenv("PORT")
	if len(port) == 0 { // if PORT is not set
		port = "8080"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		fmt.Println(err)
	}
}
