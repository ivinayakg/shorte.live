package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/ivinayakg/shorte.live/api/helpers"
)

func SystemAvailable(w http.ResponseWriter, r *http.Request) {
	result := true
	if helpers.SystemUnderMaintenance() {
		result = false
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "success", "available": result})
}
