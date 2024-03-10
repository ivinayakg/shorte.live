package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/constants"
	"github.com/ivinayakg/shorte.live/api/controllers"
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
	router.HandleFunc("/system/available", controllers.SystemAvailable).Methods("GET")
	router.NotFoundHandler = http.HandlerFunc(controllers.NotFound)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	}).Methods("GET", "POST", "PATCH", "DELETE")
}

func createRouter() *http.Handler {
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
	helpers.SetupTracker(time.Second*10, 200, 0)
	r.Use(middleware.LogMW)
	r.Use(middleware.OriginHandler)

	go helpers.Tracker.StartFlush()

	setupRoutes(r)
	routerProtected := corsHandler.Handler(r)
	return &routerProtected
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	PORT := os.Getenv("PORT")
	helpers.ENV = os.Getenv("ENV")

	go func() {
		if os.Getenv("ENV") == string(constants.Prod) {
			return
		}
		router := createRouter()
		fmt.Println("Starting the server on port " + "5100")
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", 5100), *router))
	}()

	router := createRouter()
	fmt.Println("Starting the server on port " + PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), *router))
}
