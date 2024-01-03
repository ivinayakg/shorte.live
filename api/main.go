package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/go/url-shortner/helpers"
	"example.com/go/url-shortner/middleware"
	"example.com/go/url-shortner/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func setupRoutes(router *mux.Router) {
	routes.UserRoutes(router.PathPrefix("/user").Subrouter())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	r := mux.NewRouter()
	helpers.CreateDBInstance()
	r.Use(middleware.LogMW)

	setupRoutes(r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", os.Getenv("PORT")), r))
}
