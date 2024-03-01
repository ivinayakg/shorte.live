package integration_tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/ivinayakg/shorte.live/api/models"
	"github.com/ivinayakg/shorte.live/api/utils"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestGoogleLogin(t *testing.T) {
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := httpClient.Get(ServerURL + "/user/sign_in_with_google")
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	assert.Equal(t, http.StatusFound, resp.StatusCode, "Expected status code to be 302")
	assert.Contains(t, resp.Header.Get("Location"), "https://accounts.google.com/o/oauth2/auth", "Expected redirect to Google OAuth URL")
}

func TestSelfUser(t *testing.T) {
	user := models.User{}

	TestDb.User.FindOne(context.Background(), bson.M{"email": "test1@gmail.com"}).Decode(&user)

	userJwt, _ := utils.CreateJWT(&user)
	authCookie := utils.CreateAuthCookie(*userJwt)

	req, err := http.NewRequest(http.MethodGet, ServerURL+"/user/self", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	req.AddCookie(authCookie)

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	body := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&body)
	fmt.Println(body)

	assert.Equal(t, resp.StatusCode, http.StatusOK, "Expected status code to be 200")
	assert.IsType(t, body["_id"], "string", "Expected _id to be a string")
	assert.IsType(t, body["token"], "string", "Expected token to be a string")
}

func TestSelfUserUnauthenticated(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, ServerURL+"/user/self", nil)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	body := map[string]interface{}{}

	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, resp.StatusCode, http.StatusForbidden, "Expected status code to be 200")
}
