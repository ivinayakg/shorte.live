package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"example.com/go/url-shortner/helpers"
	"example.com/go/url-shortner/middleware"
	"example.com/go/url-shortner/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func setupRoutes(router *mux.Router) {
	routes.UserRoutes(router.PathPrefix("/user").Subrouter())
	routes.URLResolveRoutes(router)
	routes.URLRoutes(router.PathPrefix("/url").Subrouter())
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	allowed_origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), " ")
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowed_origins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r := mux.NewRouter()
	helpers.CreateDBInstance()
	helpers.RedisSetup()
	r.Use(middleware.LogMW)

	setupRoutes(r)
	routerProtected := corsHandler.Handler(r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", os.Getenv("PORT")), routerProtected))
}
