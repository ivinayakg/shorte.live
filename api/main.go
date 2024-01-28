package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/middleware"
	"github.com/ivinayakg/shorte.live/api/routes"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func setupRoutes(router *mux.Router) {
	routes.UserRoutes(router.PathPrefix("/user").Subrouter())
	routes.URLResolveRoutes(router)
	routes.URLRoutes(router.PathPrefix("/url").Subrouter())
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	PORT := os.Getenv("PORT")

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

	fmt.Println("Starting the server on port " + PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), routerProtected))
}
