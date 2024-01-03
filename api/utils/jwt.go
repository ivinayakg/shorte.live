package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"example.com/go/url-shortner/models"
	"github.com/golang-jwt/jwt"
)

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))
var expiry = os.Getenv("JWT_EXPIRY")

func CreateJWT(user *models.User) (*string, error) {
	expiryTotal, err := strconv.Atoi(expiry)
	if err != nil {
		fmt.Println("Error:", err)
		expiryTotal = 21600
	}

	if expiryTotal <= 0 {
		return nil, fmt.Errorf("invalid expiry value: %s", expiry)
	}

	expirationTime := time.Now().Add(time.Duration(expiryTotal) * time.Second)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"email":  user.Email,
		"exp":    expirationTime.Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

func VerifyJwt(tokenString string) (*map[string]string, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("invalid token value: %s", tokenString)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &map[string]string{"userId": fmt.Sprint(claims["userId"]), "email": fmt.Sprint(claims["email"])}, nil
	} else {
		return nil, fmt.Errorf("failed to extract claims from token")
	}
}
