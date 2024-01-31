package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/middleware"
	"github.com/ivinayakg/shorte.live/api/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShortenURLRequest struct {
	Destination string `json:"destination"`
	CustomShort string `json:"short"`
	Expiry      int64  `json:"expiry"`
}

type ShortenURLReponse struct {
	Destination string `json:"destination"`
	CustomShort string `json:"short"`
	Expiry      int64  `json:"expiry"`
}

type UpdateURLRequest struct {
	CustomShort string `json:"short"`
	Expiry      int64  `json:"expiry"`
	Destination string `json:"destination"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middleware.UserAuthKey).(*models.User)
	body := new(ShortenURLRequest)
	preoccupiedShorts := []string{"url", "user"}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	info, err := helpers.RateLimit(r, userData.ID.Hex(), nil)
	if err != nil {
		helpers.SendJSONError(w, http.StatusTooManyRequests, fmt.Errorf("you have exhausted your quota for %v, %v to retry again", "Shorten URL", helpers.TimeRemaining(info)).Error())
		return
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.Destination) {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid url").Error())
		return
	}

	// check for domain error
	if !helpers.RemoverDomainError(body.Destination) {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("invalid url").Error())
		return
	}

	if body.CustomShort != "" && helpers.ContainsString(&preoccupiedShorts, &body.CustomShort) && helpers.ValidShortString(&body.CustomShort) {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("can't use this short").Error())
		return
	}

	// enforce https, SSL
	body.Destination = helpers.EnforceHTTP(body.Destination)

	if body.Expiry < helpers.LowestUnixTime() {
		body.Expiry = time.Now().Add(time.Hour * 48).Unix()
	}

	shortedURL, err := models.CreateURL(userData, body.CustomShort, body.Destination, body.Expiry)
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := ShortenURLReponse{
		Destination: shortedURL.Destination,
		CustomShort: shortedURL.Short,
		Expiry:      int64(shortedURL.Expiry),
	}

	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
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
		if !currentTime.After(time.Unix(int64(url.Expiry), 0)) {
			urlExpiredOrNotFound = false
		}
	} else {
		url, err = models.GetURL(urlShort, "")
		if err != nil && err != mongo.ErrNoDocuments {
			helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err != mongo.ErrNoDocuments && !currentTime.After(time.Unix(int64(url.Expiry), 0)) {
			urlExpiredOrNotFound = false
			go helpers.Redis.SetJSON(urlShort, url, time.Until(time.Unix(int64(url.Expiry), 0)))
		}
	}

	if urlExpiredOrNotFound || url == nil {
		notFoundUrl := os.Getenv("FRONTEND_URL") + "/not-found/redirect"
		http.Redirect(w, r, notFoundUrl, http.StatusMovedPermanently)
		return
	}

	// go func(urlId string) {
	// 	currTime := time.Now()
	// 	models.UpdateUserURLVisited(urlId, currTime)
	// }(url.ID.Hex())

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
}

func GetUserURL(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middleware.UserAuthKey).(*models.User)

	urls, err := models.GetUserURL(userData.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	helpers.SetHeaders("GET", w, http.StatusOK)

	if len(urls) == 0 {
		json.NewEncoder(w).Encode([]map[string]string{})
		return
	}

	json.NewEncoder(w).Encode(urls)
}

func UpdateUrl(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middleware.UserAuthKey).(*models.User)
	vars := mux.Vars(r)
	urlId := vars["id"]
	reqData := new(UpdateURLRequest)

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	url, err := models.GetURL("", urlId)
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if reqData.CustomShort == "" {
		reqData.CustomShort = url.Short
	}
	if reqData.Destination == "" {
		reqData.Destination = url.Destination
	}

	var expiry = models.UnixTime(reqData.Expiry)
	if reqData.Expiry < helpers.LowestUnixTime() {
		expiry = url.Expiry
	}

	if err := models.UpdateUserURL(userData.ID, urlId, reqData.CustomShort, reqData.Destination, expiry); err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	go helpers.Redis.Client.Del(context.Background(), url.Short)

	helpers.SetHeaders("PATCH", w, http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "successfully updated"})
}

func DeleteUrl(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middleware.UserAuthKey).(*models.User)
	vars := mux.Vars(r)
	urlId := vars["id"]

	url, err := models.GetURL("", urlId)
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := models.DeleteURL(userData.ID, urlId); err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	go helpers.Redis.Client.Del(context.Background(), url.Short)

	helpers.SetHeaders("DELETE", w, http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "successfully deleted"})
}
