package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/go/url-shortner/middleware"
	"example.com/go/url-shortner/models"
	"example.com/go/url-shortner/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
)

func setupRoutes(router *mux.Router, db *bun.DB) {
	routes.UserRoutes(router.PathPrefix("/user").Subrouter())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	db, _ := models.ConnectToDB()
	r := mux.NewRouter()
	r.Use(middleware.LogMW)

	setupRoutes(r, db)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", os.Getenv("PORT")), r))
}
