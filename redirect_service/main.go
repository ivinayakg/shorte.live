package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/controllers"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/middleware"
	"github.com/joho/godotenv"
)

var UI_URL string

func RedirectHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, UI_URL, http.StatusSeeOther)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	PORT := os.Getenv("PORT")
	UI_URL = os.Getenv("UI_URL")

	r := mux.NewRouter()
	helpers.CreateDBInstance()
	helpers.RedisSetup()
	helpers.SetupTracker(time.Second*10, 200, 0)
	r.Use(middleware.LogMW)

	go helpers.Tracker.StartFlush()

	r.HandleFunc("/", RedirectHome).Methods("GET", "POST", "PUT", "PATCH", "DELETE")
	r.HandleFunc("/{short}", controllers.ResolveURL).Methods("GET")

	fmt.Println("Starting the server on port " + PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), r))
}
