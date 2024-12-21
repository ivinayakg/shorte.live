package models

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ivinayakg/shorte.live/api/database"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func CreateURL(user *database.User, short string, destination string, expiry int64, temp bool) (*database.URL, error) {
	url := new(database.URL)

	if short != "" {
		err := database.CurrentDb.Url.FindOne(context.TODO(), bson.M{"short": short}).Decode(&url)
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

	if user != nil {
		url.User = user.ID
	}
	url.Short = short
	url.Destination = destination
	url.Expiry = database.UnixTime(expiry)
	// url.UpdateAt = database.UnixTime(time.Now().Unix())
	url.CreatedAt = database.UnixTime(time.Now().Unix())
	url.ID = primitive.NilObjectID
	url.Temporary = temp
	ctx := context.TODO()

	res, err := database.CurrentDb.Url.InsertOne(ctx, url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	url.ID = res.InsertedID.(primitive.ObjectID)
	fmt.Printf("URL created with id %v\n", url.ID)

	url.Short = helpers.BuildUrl("/" + url.Short)
	url.UserDoc = user

	return url, nil
}

func GetURL(short string, id string) (*database.URL, error) {
	url := new(database.URL)

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

	err := database.CurrentDb.Url.FindOne(ctx, urlFilter).Decode(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("URL found with id %v\n", url.ID)
	return url, nil
}

func GetUserURL(userId primitive.ObjectID) ([]*database.URL, error) {
	ctx := context.TODO()
	urlFilter := bson.M{"user": userId}

	curr, err := database.CurrentDb.Url.Find(ctx, urlFilter)
	if err != nil {
		return nil, err
	}
	defer curr.Close(context.TODO())

	var results []*database.URL
	for curr.Next(context.TODO()) {
		var result database.URL
		e := curr.Decode(&result)
		if e != nil {
			fmt.Println(err)
		}
		result.Short = helpers.BuildUrl("/" + result.Short)
		results = append(results, &result)
	}

	if err := curr.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func UpdateUserURL(userId primitive.ObjectID, urlId string, newShort string, destination string, expiry database.UnixTime) error {
	urlObjectId, err := primitive.ObjectIDFromHex(urlId)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ctx := context.TODO()
	urlFilter := bson.M{"user": userId, "_id": urlObjectId}

	// Check if the document exists
	count, err := database.CurrentDb.Url.CountDocuments(ctx, urlFilter)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if count == 0 {
		fmt.Println("URL Document not found")
		return errors.New("URL document not found")
	}

	updateData := bson.M{"$set": bson.M{"short": newShort, "destination": destination, "expiry": expiry}}

	res, err := database.CurrentDb.Url.UpdateOne(ctx, urlFilter, updateData)
	if err != nil {
		if writeError, ok := err.(mongo.WriteException); ok && writeError.WriteErrors[0].Code == 11000 {
			fmt.Println(err)
			return fmt.Errorf("URL custom short is already in user")
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

	res, err := database.CurrentDb.Url.UpdateOne(ctx, urlFilter, updateData)
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

	res, err := database.CurrentDb.Url.DeleteOne(ctx, urlFilter)
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
