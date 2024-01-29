package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ivinayakg/shorte.live/api/helpers"
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

	url.Short = os.Getenv("DOMAIN") + "/" + url.Short

	return &url, nil
}

func GetURL(short string, id string) (*URLDoc, error) {
	url := new(URLDoc)

	ctx := context.TODO()

	var urlFilter bson.M
	if id == "" {
		urlFilter = bson.M{"short": short}
	} else {
		urlObjectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		urlFilter = bson.M{"_id": urlObjectId}
	}

	err := helpers.CurrentDb.Url.FindOne(ctx, urlFilter).Decode(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("URL found with id %v\n", url.ID)
	return url, nil
}

func GetUserURL(userId primitive.ObjectID) ([]*URLDoc, error) {
	ctx := context.TODO()
	urlFilter := bson.M{"user": userId}

	curr, err := helpers.CurrentDb.Url.Find(ctx, urlFilter)
	if err != nil {
		return nil, err
	}
	defer curr.Close(context.TODO())

	var results []*URLDoc
	for curr.Next(context.TODO()) {
		var result URLDoc
		e := curr.Decode(&result)
		if e != nil {
			fmt.Println(err)
		}
		result.Short = os.Getenv("SHORTED_URL_DOMAIN") + "/" + result.Short
		results = append(results, &result)
	}

	if err := curr.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func UpdateUserURL(userId primitive.ObjectID, urlId string, newShort string, destination string, expiry time.Time) error {
	urlObjectId, err := primitive.ObjectIDFromHex(urlId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// URL with newShort already exists
	var urlDoc = new(URLDoc)
	err = helpers.CurrentDb.Url.FindOne(context.TODO(), bson.M{"short": newShort}).Decode(&urlDoc)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println(err)
		return err
	}
	if urlDoc.ID != primitive.NilObjectID && urlDoc.ID.Hex() != urlId {
		return fmt.Errorf("URL custom short is already in user")
	}

	ctx := context.TODO()
	urlFilter := bson.M{"user": userId, "_id": urlObjectId}
	updateData := bson.M{"$set": bson.M{"short": newShort, "destination": destination, "expiry": expiry}}

	res, err := helpers.CurrentDb.Url.UpdateOne(ctx, urlFilter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("URL Document not found")
		} else {
			fmt.Println(err)
		}
		return err
	} else {
		fmt.Printf("Update document successfully URL: %+v\n", res.UpsertedID)
	}

	return nil
}

func UpdateUserURLVisited(urlId string, visited time.Time) error {
	urlObjectId, err := primitive.ObjectIDFromHex(urlId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx := context.TODO()
	urlFilter := bson.M{"_id": urlObjectId}
	updateData := bson.M{"$set": bson.M{"lastvisited": visited}}

	res, err := helpers.CurrentDb.Url.UpdateOne(ctx, urlFilter, updateData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("URL Document not found")
		} else {
			fmt.Println(err)
		}
		return err
	} else {
		fmt.Printf("Update document successfully URL: %+v\n", res.UpsertedID)
	}

	return nil
}

func DeleteURL(userId primitive.ObjectID, urlId string) error {
	urlObjectId, err := primitive.ObjectIDFromHex(urlId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx := context.TODO()
	urlFilter := bson.M{"user": userId, "_id": urlObjectId}

	res, err := helpers.CurrentDb.Url.DeleteOne(ctx, urlFilter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Println("URL Document not found")
		} else {
			fmt.Println(err)
		}
		return err
	} else {
		fmt.Printf("Deleted document successfully URL: %+v. Total documents deleted - %d\n", urlId, res.DeletedCount)
	}

	return nil
}
