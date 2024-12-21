package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/avct/uasurfer"
	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/constants"
	"github.com/ivinayakg/shorte.live/api/database"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/middleware"
	"github.com/ivinayakg/shorte.live/api/models"
	"github.com/ivinayakg/shorte.live/api/timescale"
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
	userData := r.Context().Value(middleware.UserAuthKey).(*database.User)
	body := new(ShortenURLRequest)

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

	if body.CustomShort != "" && (helpers.ContainsString(&constants.PreoccupiedShorts, &body.CustomShort) || helpers.NotValidShortString(&body.CustomShort)) {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("can't use this short").Error())
		return
	}

	// enforce https, SSL
	body.Destination = helpers.EnforceHTTP(body.Destination)

	if body.Expiry < helpers.LowestUnixTime() {
		body.Expiry = time.Now().Add(time.Hour * 48).Unix()
	}

	shortedURL, err := models.CreateURL(userData, body.CustomShort, body.Destination, body.Expiry, false)
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
	url := &database.URL{}
	urlExpiredOrNotFound := true
	var err error

	revalidateCache, err := strconv.ParseBool(r.URL.Query().Get("revalidate"))
	if err != nil {
		fmt.Println("Error:", err)
		revalidateCache = false
	}

	systemNotAvailable := helpers.SystemUnderMaintenance(revalidateCache)
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
		notFoundUrl := os.Getenv("UI_NOT_FOUND_URL")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		http.Redirect(w, r, notFoundUrl, http.StatusTemporaryRedirect)
		return
	}

	go func(r *http.Request, url database.URL) {
		if url.Temporary {
			return
		}

		userAgent := r.Header.Get("User-Agent")
		ua := uasurfer.Parse(userAgent)

		device := ua.DeviceType.String()
		ip := strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0]
		os := ua.OS.Name.String()
		referrer := r.Header.Get("Referer")
		urlId := url.ID.Hex()

		if referrer == "" {
			referrer = "direct"
		}

		timestamp := time.Now().Unix()

		if helpers.ENV != string(constants.Prod) {
			helpers.Tracker.CaptureRedirectEvent(device, ip, os, referrer, urlId, timestamp)
		}

	}(r, *url)

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, url.Destination, http.StatusMovedPermanently)
}

func GetUserURL(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middleware.UserAuthKey).(*database.User)

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
	userData := r.Context().Value(middleware.UserAuthKey).(*database.User)
	vars := mux.Vars(r)
	urlId := vars["id"]
	reqData := new(UpdateURLRequest)

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	url, err := models.GetURL("", urlId)
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if reqData.CustomShort == "" {
		reqData.CustomShort = url.Short
	}
	if reqData.Destination == "" {
		reqData.Destination = url.Destination
	}

	var expiry = database.UnixTime(reqData.Expiry)
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
	userData := r.Context().Value(middleware.UserAuthKey).(*database.User)
	vars := mux.Vars(r)
	urlId := vars["id"]

	url, err := models.GetURL("", urlId)
	if err != nil {
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

func GetURLStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlId := vars["id"]
	userData := r.Context().Value(middleware.UserAuthKey).(*database.User)

	startTimeStr := r.URL.Query().Get("start")
	endTimeStr := r.URL.Query().Get("end")

	var startTime, endTime int64

	if startTimeStr == "" {
		startTime = time.Now().Add(-time.Hour * 24).Unix()
	} else {
		startTime, _ = strconv.ParseInt(startTimeStr, 10, 64)
		if startTime < time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Unix() || startTime > time.Now().Unix() {
			helpers.SendJSONError(w, http.StatusBadRequest, "invalid start time")
			return
		}
	}

	if endTimeStr == "" {
		endTime = time.Now().Unix()
	} else {
		endTime, _ = strconv.ParseInt(endTimeStr, 10, 64)
		if endTime > time.Now().Unix() {
			helpers.SendJSONError(w, http.StatusBadRequest, "invalid end time")
			return
		}
	}

	url, err := models.GetURL("", urlId)
	if err != nil || url == nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	if url.User.Hex() != userData.ID.Hex() {
		helpers.SendJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	stats, err := timescale.QueryClickEvents(urlId, startTime, endTime)
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SetHeaders("GET", w, http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "successfully updated", "data": stats})
}

func ShortenURLTemp(w http.ResponseWriter, r *http.Request) {
	body := new(ShortenURLRequest)

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	rateConfig := &helpers.URLLimit{Value: 5, Expiry: 1440}

	info, err := helpers.RateLimit(r, "", rateConfig)
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

	// enforce https, SSL
	body.Destination = helpers.EnforceHTTP(body.Destination)
	body.Expiry = time.Now().Add(time.Hour * 48).Unix()

	shortedURL, err := models.CreateURL(nil, "", body.Destination, body.Expiry, true)
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
