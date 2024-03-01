package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ivinayakg/shorte.live/api/constants"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/models"
	"github.com/ivinayakg/shorte.live/api/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type userAuth string

const UserAuthKey userAuth = "User"

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		cookie := utils.GetCookie(r)
		if cookie != nil {
			token = cookie.Value
		} else if helpers.ENV != string(constants.Prod) {
			tokenHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
			if len(tokenHeader) < 2 {
				errMsg := "Authentication error!, Provide valid auth token"
				helpers.SendJSONError(w, http.StatusForbidden, errMsg)
				log.Println(errMsg)
				return
			}
			token = tokenHeader[1]
		} else {
			errMsg := "Authentication error!, login first"
			helpers.SendJSONError(w, http.StatusForbidden, errMsg)
			log.Println(errMsg)
			return
		}

		systemNotAvailable := helpers.SystemUnderMaintenance(false)
		if systemNotAvailable {
			error := fmt.Errorf("system is under maintenance")
			helpers.SendJSONError(w, http.StatusServiceUnavailable, error.Error())
			return
		}

		verifyUserData, err := utils.VerifyJwt(token)
		if err != nil {
			errMsg := err.Error()
			helpers.SendJSONError(w, http.StatusForbidden, errMsg)
			log.Println(errMsg)
			return
		}

		user, err := models.GetUser((*verifyUserData)["email"])
		if err != nil {
			errMsg := err.Error()
			if err != mongo.ErrNoDocuments {
				errMsg = "Authentication error!"
			}
			helpers.SendJSONError(w, http.StatusForbidden, errMsg)
			log.Println(errMsg)
			return
		}

		if helpers.ENV != string(constants.Prod) {
			user.Token = token
		}

		c := context.WithValue(r.Context(), UserAuthKey, user)
		next.ServeHTTP(w, r.WithContext(c))
	})
}
