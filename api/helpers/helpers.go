package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

var methodChoices = map[string]string{
	"get":   "GET",
	"post":  "POST",
	"patch": "PATCH",
	"del":   "DELETE",
}

func EnforceHTTP(url string) string {
	if url[:4] != "http" {
		return "https://" + url
	}
	return url
}

func RemoverDomainError(url string) bool {
	if url == os.Getenv("DOMAIN") {
		return false
	}

	newUrl := strings.Replace(url, "http://", "", 1)
	newUrl = strings.Replace(newUrl, "https://", "", 1)
	newUrl = strings.Replace(newUrl, "www.", "", 1)
	newUrl = strings.Split(newUrl, "/")[0]

	if newUrl == os.Getenv("DOMAIN") {
		return false
	}

	return true
}
func SetHeaders(type_ string, w http.ResponseWriter, status int) {
	method := methodChoices[type_]
	if method == "" {
		method = "GET"
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if method != "GET" {
		w.Header().Set("Access-Control-Allow-Methods", method)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	}
}

func SendJSONError(w http.ResponseWriter, statusCode int, errorMessage string) {
	errorResponse := ErrorResponse{Error: errorMessage}
	w.Header().Set("content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func ContainsString(arr *[]string, target *string) bool {
	for _, s := range *arr {
		if strings.Contains(s, *target) {
			return true
		}
	}
	return false
}

func GetUserIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")

	// If X-Forwarded-For header is empty (not behind a proxy), get the IP from RemoteAddr
	if ip == "" {
		ip = r.RemoteAddr
	}

	return ip
}

func TimeRemaining(duration time.Duration) string {
	if duration <= 0 {
		return "Time has expired"
	}

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return fmt.Sprintf("Time remaining = %02dd %02dh %02dm %02ds", days, hours, minutes, seconds)
}

func NotValidShortString(short *string) bool {
	re := regexp.MustCompile(`[/@&?#.]+`)
	return re.MatchString(*short)
}

func LowestUnixTime() int64 {
	return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()
}
