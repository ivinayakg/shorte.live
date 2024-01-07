package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"example.com/go/url-shortner/helpers"
	"example.com/go/url-shortner/middleware"
	"example.com/go/url-shortner/models"
	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShortenURLRequest struct {
	Destination string `json:"destination"`
	CustomShort string `json:"short"`
	Expiry      int32  `json:"expiry"`
}

type ShortenURLReponse struct {
	Destination string    `json:"destination"`
	CustomShort string    `json:"short"`
	Expiry      time.Time `json:"expiry"`
}

type UpdateURLRequest struct {
	CustomShort string `json:"short"`
	Expiry      string `json:"expiry"`
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

	//*implement rate limiting
	// r2 := database.CreateClient(1)
	// defer r2.Close()
	// val, err := r2.Get(database.Ctx, c.IP()).Result()
	// if err == redis.Nil {
	// 	_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	// } else {
	// 	// val, _ = r2.Get(database.Ctx, c.IP()).Result()
	// 	valInt, _ := strconv.Atoi(val)
	// 	if valInt <= 0 {
	// 		limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
	// 		return {"error": "rate limit exceeded", "rate_limit_rest": limit / time.Nanosecond / time.Minute}
	// 	}
	// }

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

	if body.CustomShort != "" && helpers.ContainsString(&preoccupiedShorts, &body.CustomShort) {
		helpers.SendJSONError(w, http.StatusBadRequest, fmt.Errorf("short already in use").Error())
		return
	}

	// enforce https, SSL
	body.Destination = helpers.EnforceHTTP(body.Destination)

	// r := database.CreateClient(0)
	// defer r.Close()
	// val, _ = r.Get(database.Ctx, id).Result()
	// if val != "" {
	// 	return "URL custom short is already in user"
	// }

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	tinyUrl, err := models.CreateURL(userData, body.CustomShort, body.Destination, body.Expiry)
	if err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	// if err != nil {
	// 	return {"error": "unable to connect to server"}
	// }
	resp := ShortenURLReponse{
		Destination: tinyUrl.Destination,
		CustomShort: tinyUrl.Short,
		Expiry:      tinyUrl.Expiry,
	}

	// r2.Decr(database.Ctx, c.IP())

	// val, _ = r2.Get(database.Ctx, c.IP()).Result()

	// resp.XRateRemaining, _ = strconv.Atoi(val)
	// ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	// resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func ResolveURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlShort := vars["short"]

	url, err := models.GetURL(urlShort, "")
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	currentTime := time.Now()
	if err == mongo.ErrNoDocuments || currentTime.After(url.Expiry) {
		notFoundUrl := os.Getenv("FRONTEND_URL") + "/not-found/redirect"
		http.Redirect(w, r, notFoundUrl, http.StatusMovedPermanently)
		return
	}

	// go func(urlId string) {
	// 	currTime := time.Now()
	// 	models.UpdateUserURLVisited(urlId, currTime)
	// }(url.ID.Hex())

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

	var expiry time.Time
	if reqData.Expiry != "" {
		parsedDatetime, err := time.Parse(time.RFC3339, reqData.Expiry)
		if err != nil {
			helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		expiry = parsedDatetime
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
	if reqData.Expiry == "" {
		expiry = url.Expiry
	}

	if err := models.UpdateUserURL(userData.ID, urlId, reqData.CustomShort, reqData.Destination, expiry); err != nil {
		helpers.SendJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	helpers.SetHeaders("PATCH", w, http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "successfully updated"})
}
