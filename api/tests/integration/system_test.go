package integration_tests

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/tests/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestURLSystemAvailability(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, ServerURL+"/system/available", nil)
	if err != nil {
		t.Fatal(err)
	}
	testhelper.PutSystemUnderMaintenance(helpers.Redis, false)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&body)

	assert.Equal(t, true, body["available"], "Expected maintenance to be true")
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Expected status code to be 200")
}

func TestURLSystemAvailabilityFail(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, ServerURL+"/system/available", nil)
	if err != nil {
		t.Fatal(err)
	}

	testhelper.PutSystemUnderMaintenance(helpers.Redis, true)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	body := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&body)

	assert.Equal(t, true, body["available"], "Expected maintenance to be true")
	assert.Equal(t, resp.StatusCode, http.StatusOK, "Expected status code to be 200")
}

func TestNotFound(t *testing.T) {
	resp, err := RedirecthttpClient.Get(ServerURL + "/" + "something/random")
	if err != nil {
		t.Fatal(err)
	}

	notFoundUrl := os.Getenv("UI_NOT_FOUND_URL")

	assert.Equal(t, resp.StatusCode, http.StatusTemporaryRedirect, "Excpected status code to be 307")
	assert.Contains(t, resp.Header.Get("Location"), notFoundUrl, "Expected redirect to not found url")
}

func TestRedirectHome(t *testing.T) {
	resp, err := RedirecthttpClient.Get(ServerURL + "/")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusSeeOther, "Excpected status code to be 303")
	assert.Contains(t, resp.Header.Get("Location"), os.Getenv("FRONTEND_URL"), "Expected redirect to homepage url")
}
