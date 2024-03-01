package middleware

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
)

var re = regexp.MustCompile(`^/[^/]+$`)

func OriginHandler(next http.Handler) http.Handler {
	var RedirectServiceUrl = os.Getenv("SHORTED_URL_DOMAIN")
	notFoundUrl := os.Getenv("UI_NOT_FOUND_URL")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Host, RedirectServiceUrl)
		if r.Host == RedirectServiceUrl && !re.MatchString(r.URL.Path) {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			http.Redirect(w, r, notFoundUrl, http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
