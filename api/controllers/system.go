package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/ivinayakg/shorte.live/api/helpers"
)

func SystemAvailable(w http.ResponseWriter, r *http.Request) {
	result := true
	if helpers.SystemUnderMaintenance(false) {
		result = false
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success", "available": result})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	notFoundUrl := os.Getenv("UI_NOT_FOUND_URL")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, notFoundUrl, http.StatusTemporaryRedirect)
}

func RedirectHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	http.Redirect(w, r, os.Getenv("FRONTEND_URL"), http.StatusSeeOther)
}
