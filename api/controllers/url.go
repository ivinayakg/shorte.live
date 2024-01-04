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
)

type ShortenURLRequest struct {
	URL         string `json:"url"`
	CustomShort string `json:"short"`
	Expiry      int32  `json:"expiry"`
}

type ShortenURLReponse struct {
	URL         string    `json:"url"`
	CustomShort string    `json:"short"`
	Expiry      time.Time `json:"expiry"`
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middleware.UserAuthKey).(*models.User)
	body := new(ShortenURLRequest)

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	if !govalidator.IsURL(body.URL) {
		http.Error(w, fmt.Errorf("invalid url").Error(), http.StatusBadRequest)
		return
	}

	// check for domain error
	if !helpers.RemoverDomainError(body.URL) {
		http.Error(w, fmt.Errorf("invalid url").Error(), http.StatusBadRequest)
		return
	}

	// enforce https, SSL
	body.URL = helpers.EnforceHTTP(body.URL)

	// r := database.CreateClient(0)
	// defer r.Close()
	// val, _ = r.Get(database.Ctx, id).Result()
	// if val != "" {
	// 	return "URL custom short is already in user"
	// }

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	tinyUrl, err := models.CreateURL(userData, body.CustomShort, body.URL, body.Expiry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	// if err != nil {
	// 	return {"error": "unable to connect to server"}
	// }
	resp := ShortenURLReponse{
		URL:         tinyUrl.Destination,
		CustomShort: "",
		Expiry:      tinyUrl.Expiry,
	}

	// r2.Decr(database.Ctx, c.IP())

	// val, _ = r2.Get(database.Ctx, c.IP()).Result()

	// resp.XRateRemaining, _ = strconv.Atoi(val)
	// ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()
	// resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + tinyUrl.Short
	helpers.SetHeaders("post", w, http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
