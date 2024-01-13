package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"example.com/go/url-shortner/helpers"
	"example.com/go/url-shortner/middleware"
	"example.com/go/url-shortner/models"
	"example.com/go/url-shortner/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func SignInWithGoogle(w http.ResponseWriter, r *http.Request) {
	tokenRequestUrl := os.Getenv("GOOGLE_OAUTH_AUTH_REQUEST_URI")
	responseCode := "code"
	clientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	redirectUri := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
	scope := "https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email"

	url := fmt.Sprintf("%v?response_type=%v&client_id=%v&redirect_uri=%v&scope=%v", tokenRequestUrl, responseCode, clientId, redirectUri, scope)

	http.Redirect(w, r, url, http.StatusFound)
}

func CallbackSignInWithGoogle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered from panic:", r)
			helpers.SendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
	}()

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form:", err)
		helpers.SendJSONError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	if r.FormValue("error") != "" || r.FormValue("code") == "" {
		http.Redirect(w, r, os.Getenv("BASE_UI_URL"), http.StatusSeeOther)
		return
	}

	accessTokenURI := os.Getenv("GOOGLE_OAUTH_TOKEN_REQUEST_URI")
	redirectURI := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")

	params := url.Values{
		"code":          {r.FormValue("code")},
		"redirect_uri":  {redirectURI},
		"client_id":     {os.Getenv("GOOGLE_OAUTH_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")},
		"grant_type":    {"authorization_code"},
	}

	resp, err := http.PostForm(accessTokenURI, params)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		helpers.SendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	var tokenData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenData); err != nil {
		fmt.Println("Error decoding token data:", err)
		helpers.SendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	userInfoURI := fmt.Sprintf("/oauth2/v1/userinfo?access_token=%s", tokenData["access_token"])

	resp, err = http.Get("https://www.googleapis.com" + userInfoURI)
	if err != nil {
		fmt.Println("Error making GET request:", err)
		helpers.SendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	defer resp.Body.Close()

	var googleProfile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&googleProfile); err != nil {
		fmt.Println("Error decoding Google profile:", err)
		helpers.SendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	fmt.Println(googleProfile)

	// Assuming User, UserSerializer, and other settings are defined elsewhere
	user, err := models.GetUser(googleProfile["email"].(string))
	if err != nil && err != mongo.ErrNoDocuments {
		helpers.SendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if user != nil {
		token, _ := utils.CreateJWT(user)
		http.Redirect(w, r, fmt.Sprintf(os.Getenv("FRONTEND_AUTH_URL")+"%v", *token), http.StatusSeeOther)
		return
	} else {
		user, err = models.CreateUser(googleProfile["email"].(string), googleProfile["name"].(string), googleProfile["picture"].(string))
		if err != nil {
			helpers.SendJSONError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		token, _ := utils.CreateJWT(user)
		http.Redirect(w, r, fmt.Sprintf(os.Getenv("FRONTEND_AUTH_URL")+"%v", *token), http.StatusSeeOther)
		return
	}

}

func SelfUser(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(middleware.UserAuthKey)
	json.NewEncoder(w).Encode(userData)
}
