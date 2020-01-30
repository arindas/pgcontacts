package main

import (
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"net/http"

	"github.com/arindas/pgcontacts/internal/app"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.AuthMiddleWare)

	port := os.Getenv("PORT")
	if len(port) == 0 { // if PORT is not set
		port = "8080"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		fmt.Println(err)
	}
}
