package models

import (
	"context"
	"fmt"
	"time"

	"example.com/go/url-shortner/helpers"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func CreateURL(user *User, short string, destination string, expiry int32) (*URL, error) {
	urlDoc := new(URLDoc)

	if short != "" {
		err := helpers.CurrentDb.Url.FindOne(context.TODO(), bson.M{"short": short}).Decode(&urlDoc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("url Document not found")
			} else {
				fmt.Println(err)
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("URL custom short is already in user")
		}
	} else {
		short = uuid.New().String()[:10]
	}

	urlDoc.User = user.ID
	urlDoc.Short = short
	urlDoc.Destination = destination
	urlDoc.Expiry = time.Now().Add(time.Duration(expiry) * 3600 * time.Second)
	urlDoc.LastVisited = time.Now()
	urlDoc.ID = primitive.NilObjectID
	ctx := context.TODO()

	res, err := helpers.CurrentDb.Url.InsertOne(ctx, urlDoc)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	url := URL{User: user, Short: short, Destination: destination, Expiry: urlDoc.Expiry, LastVisited: urlDoc.LastVisited}
	url.ID = res.InsertedID.(primitive.ObjectID)
	fmt.Printf("URL created with id %v\n", url.ID)

	return &url, nil
}

func GetURL(short string) (*URLDoc, error) {
	url := new(URLDoc)

	ctx := context.TODO()
	urlFilter := bson.M{"short": short}

	err := helpers.CurrentDb.Url.FindOne(ctx, urlFilter).Decode(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("URL found with id %v\n", url.ID)
	return url, nil
}
