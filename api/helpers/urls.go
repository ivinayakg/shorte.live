package helpers

import (
	"os"
	"strings"
)

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

func BuildUrl(url string) string {
	env := os.Getenv("ENV")
	if env == "development" {
		return "http://" + os.Getenv("SHORTED_URL_DOMAIN") + url
	}
	return "https://" + os.Getenv("SHORTED_URL_DOMAIN") + url
}
