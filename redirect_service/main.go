package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/middleware"
	"github.com/ivinayakg/shorte.live/api/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var UI_URL string

func RedirectHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, UI_URL, http.StatusSeeOther)
}

func ResolveURL(w http.ResponseWriter, r *http.Request) {
	url := &models.URL{}
	urlExpiredOrNotFound := true
	var err error

	systemNotAvailable := helpers.SystemUnderMaintenance()
	if systemNotAvailable {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(w, r, os.Getenv("FRONTEND_URL_MAINTENANCE"), http.StatusMovedPermanently)
		return
	}

	defaultLimit, found := helpers.GetRateConfig(false).Limit["dynamic"]
	if !found {
		defaultLimit = &helpers.URLLimit{Value: 100, Expiry: 30}
	}
	info, err := helpers.RateLimit(r, "", defaultLimit)
	if err != nil {
		helpers.SendJSONError(w, http.StatusTooManyRequests, fmt.Errorf("you have exhausted your quota for %v, %v to retry again", "Resolve URL", helpers.TimeRemaining(info)).Error())
		return
	}

	vars := mux.Vars(r)
	urlShort := vars["short"]
	currentTime := time.Now()

	revalidateCache, err := strconv.ParseBool(r.URL.Query().Get("revalidate"))
	if err != nil {
		fmt.Println("Error:", err)
		revalidateCache = false
	}

	err = helpers.Redis.GetJSON(urlShort, url)
	if err != nil {
		fmt.Println(err)
	}

	if url.ID != primitive.NilObjectID && !revalidateCache {
		if !currentTime.After(url.Expiry) {
			urlExpiredOrNotFound = false
		}
	} else {
		url, err = models.GetURL(urlShort, "")
		if err != nil && err != mongo.ErrNoDocuments {
			helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err != mongo.ErrNoDocuments && !currentTime.After(url.Expiry) {
			urlExpiredOrNotFound = false
			go helpers.Redis.SetJSON(urlShort, url, time.Until(url.Expiry))
		}
	}

	if urlExpiredOrNotFound || url == nil {
		notFoundUrl := os.Getenv("UI_NOT_FOUND_URL")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(w, r, notFoundUrl, http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	PORT := os.Getenv("PORT")
	UI_URL = os.Getenv("UI_URL")

	r := mux.NewRouter()

	helpers.RedisSetup()
	helpers.CreateDBInstance()
	r.Use(middleware.LogMW)

	r.HandleFunc("/", RedirectHome).Methods("GET", "POST", "PUT", "PATCH", "DELETE")
	r.HandleFunc("/{short}", ResolveURL).Methods("GET")

	fmt.Println("Starting the server on port " + PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), r))
}
