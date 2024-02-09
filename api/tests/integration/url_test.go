package integration_tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ivinayakg/shorte.live/api/models"
	"github.com/ivinayakg/shorte.live/api/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// for requests with redirects
var RedirecthttpClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}
var HttpClient = &http.Client{}

// resolve url
func TestURLResolve(t *testing.T) {
	resp, err := RedirecthttpClient.Get(ServerURL + "/" + URLFixture.Short)
	if err != nil {
		t.Fatal(err)
	}

	destinationURL := URLFixture.Destination

	assert.Equal(t, resp.StatusCode, http.StatusMovedPermanently, "Excpected status code to be 301")
	assert.Contains(t, resp.Header.Get("Location"), destinationURL, "Expected redirect to destination url")
}

func TestURLResolveNotFound(t *testing.T) {
	resp, err := RedirecthttpClient.Get(ServerURL + "/random")
	if err != nil {
		t.Fatal(err)
	}

	notFoundurl := os.Getenv("UI_NOT_FOUND_URL")

	assert.Equal(t, resp.StatusCode, http.StatusTemporaryRedirect, "Excpected status code to be 307")
	assert.Contains(t, resp.Header.Get("Location"), notFoundurl, "Expected redirect to url-not-found page")
}

func TestURLResolveExpired(t *testing.T) {
	resp, err := RedirecthttpClient.Get(ServerURL + "/" + ExpiredURLFixture.Short)
	if err != nil {
		t.Fatal(err)
	}

	notFoundurl := os.Getenv("UI_NOT_FOUND_URL")

	assert.Equal(t, resp.StatusCode, http.StatusTemporaryRedirect, "Excpected status code to be 302")
	assert.Contains(t, resp.Header.Get("Location"), notFoundurl, "Expected redirect to url-not-found page")
}

// test get user urls
func TestGetUserURLs(t *testing.T) {
	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodGet, ServerURL+"/url/all", nil)

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := []map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusOK, "Excpected status code to be 200")
	assert.Equal(t, len(respBody), 2, "Expected 2 urls")
}

// test create short url
func TestCreateShortedUrl(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"destination": "http://google.com",
		"short":       "new-short",
		"expiry":      time.Now().Add(time.Hour * 5).Unix(),
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/url", bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusCreated, "Excpected status code to be 201")
	assert.Equal(t, respBody["destination"], payloadData["destination"], "Expected destination to be http://google.com")
	assert.Contains(t, respBody["short"], payloadData["short"], "Expected short to be new-short")
}

func TestCreateShortedUrlWithInvalidUrl(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"destination": "not-a-url",
		"short":       "new-short",
		"expiry":      time.Now().Add(time.Hour * 5).Unix(),
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/url", bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "invalid url", "Expected destination to be invalid url")
}

func TestCreateShortedUrlWithInvalidUrl2(t *testing.T) {
	// destination as the deployed url
	destination := os.Getenv("DOMAIN")

	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"destination": destination,
		"short":       "new-short",
		"expiry":      time.Now().Add(time.Hour * 5).Unix(),
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/url", bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "invalid url", "Expected destination to be invalid url")
}

func TestCreateShortedUrlWithPreoccupiedShort(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"destination": "https://www.google.com",
		"short":       "user",
		"expiry":      time.Now().Add(time.Hour * 5).Unix(),
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/url", bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "can't use this short", "Expected short to be a preoccupied short")
}

func TestCreateShortedUrlWithInvalidShort(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"destination": "https://www.google.com",
		"short":       "hell@",
		"expiry":      time.Now().Add(time.Hour * 5).Unix(),
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/url", bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "can't use this short", "Expected short to be invalid")
}

func TestCreateShortedUrlWithDuplicateShort(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"destination": "https://www.google.com",
		"short":       URLFixture.Short,
		"expiry":      time.Now().Add(time.Hour * 5).Unix(),
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPost, ServerURL+"/url", bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "URL custom short is already in user", "Expected short to be already in use")
}

// update url
func TestUpdateURL(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"short": "new-short-update",
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPatch, ServerURL+"/url/"+URLFixture.ID.Hex(), bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusNoContent, "Excpected status code to be 204")

	url, _ := models.GetURL("new-short-update", "")
	assert.Contains(t, url.Short, payloadData["short"], "Expected short to be new-short")
}

func TestUpdateURLInvalidId(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"short": "new-short-update",
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodPatch, ServerURL+"/url/"+primitive.NewObjectID().Hex(), bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "mongo: no documents in result", "Expected error to be no documents in result")
}

func TestUpdateURLUnauthorized(t *testing.T) {
	// Data for the payload as a map
	payloadData := map[string]interface{}{
		"short": "new-short-update",
	}

	// Marshal the map into a JSON string
	payloadJSON, _ := json.Marshal(payloadData)

	userJwt, _ := utils.CreateJWT(&UserFixture2)

	req, _ := http.NewRequest(http.MethodPatch, ServerURL+"/url/"+URLFixture.ID.Hex(), bytes.NewBuffer(payloadJSON))

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	respBody := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "URL document not found", "Expected error to be no documents in result for different user")
}

// delete url
func TestDeleteURL(t *testing.T) {
	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodDelete, ServerURL+"/url/"+ExpiredURLFixture.ID.Hex(), nil)

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusNoContent, "Excpected status code to be 204")
}

func TestDeleteURLInvalidId(t *testing.T) {
	userJwt, _ := utils.CreateJWT(&UserFixture1)

	req, _ := http.NewRequest(http.MethodDelete, ServerURL+"/url/"+primitive.NewObjectID().Hex(), nil)

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "mongo: no documents in result", "Expected error to be no documents in result")
}

func TestDeleteURLUnauthorized(t *testing.T) {
	userJwt, _ := utils.CreateJWT(&UserFixture2)
	req, _ := http.NewRequest(http.MethodDelete, ServerURL+"/url/"+primitive.NewObjectID().Hex(), nil)

	req.Header.Add("Authorization", "Bearer "+*userJwt)

	resp, err := HttpClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	var respBody map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&respBody)

	assert.Equal(t, resp.StatusCode, http.StatusBadRequest, "Excpected status code to be 400")
	assert.Equal(t, respBody["error"], "mongo: no documents in result", "Expected error to be no documents in result for different user")
}
