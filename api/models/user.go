package models

import (
	"context"
	"fmt"
	"time"

	"github.com/ivinayakg/shorte.live/api/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func CreateUser(email string, name string, picture string) (*User, error) {
	createdAt := time.Now().In(time.UTC)
	user := User{Name: name, Email: email, Picture: picture, CreatedAt: createdAt}
	ctx := context.TODO()

	res, err := helpers.CurrentDb.User.InsertOne(ctx, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)
	fmt.Printf("User created with id %v\n", user.ID)

	return &user, nil
}

func GetUser(email string) (*User, error) {
	user := new(User)

	ctx := context.TODO()
	userFilter := bson.M{"email": email}

	err := helpers.CurrentDb.User.FindOne(ctx, userFilter).Decode(user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("User found with id %v\n", user.ID)
	return user, nil
}
