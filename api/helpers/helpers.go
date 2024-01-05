package helpers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
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
