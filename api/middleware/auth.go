package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"example.com/go/url-shortner/models"
	"example.com/go/url-shortner/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type userAuth string

const UserAuthKey userAuth = "User"

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(tokenHeader) < 2 {
			errMsg := "Authentication error!, Provide valid auth token"
			http.Error(w, errMsg, http.StatusForbidden)
			log.Println(errMsg)
			return
		}
		token := tokenHeader[1]

		verifyUserData, err := utils.VerifyJwt(token)
		if err != nil {
			errMsg := "Authentication error!"
			http.Error(w, errMsg, http.StatusForbidden)
			log.Println(errMsg)
			return
		}

		user, err := models.GetUser((*verifyUserData)["email"])
		if err != nil {
			if err != mongo.ErrNoDocuments {
				errMsg := "Authentication error!"
				http.Error(w, errMsg, http.StatusForbidden)
				log.Println(errMsg)
				return
			}
		}

		user.Token = token

		c := context.WithValue(r.Context(), UserAuthKey, user)
		next.ServeHTTP(w, r.WithContext(c))
	})
}
