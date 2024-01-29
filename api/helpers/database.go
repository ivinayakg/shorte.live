package helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	User    *mongo.Collection
	Url     *mongo.Collection
	Tracker *mongo.Collection
	Config  *mongo.Collection
}

type DBIndexName string

const urlShortIndexName DBIndexName = "url_short_index_1"

func doesIndexExist(ctx context.Context, collection *mongo.Collection, indexName string) (bool, error) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)

	var indexDoc bson.M
	for cursor.Next(ctx) {
		if err := cursor.Decode(&indexDoc); err != nil {
			return false, err
		}

		// Check if the index name matches
		if name, ok := indexDoc["name"].(string); ok && name == indexName {
			return true, nil
		}
	}

	return false, nil
}

var CurrentDb *DB

func CreateDBInstance() {
	connectionString := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	userCollName := os.Getenv("DB_USER_COLLECTION_NAME")
	urlCollName := os.Getenv("DB_URL_COLLECTION_NAME")
	trackerCollName := os.Getenv("DB_TRACKER_COLLECTION_NAME")
	configCollName := os.Getenv("DB_CONFIG_COLLECTION_NAME")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
		return
	}

	fmt.Println("Connected to MongoDB")

	userCollection := client.Database(dbName).Collection(userCollName)
	urlCollection := client.Database(dbName).Collection(urlCollName)
	trackerCollection := client.Database(dbName).Collection(trackerCollName)
	configCollection := client.Database(dbName).Collection(configCollName)

	urlShortIndex, err := doesIndexExist(context.Background(), urlCollection, string(urlShortIndexName))
	if err != nil {
		log.Fatal(err)
	}

	if !urlShortIndex {
		// Create the index
		indexModel := mongo.IndexModel{
			Keys:    bson.M{"short": 1},
			Options: options.Index().SetUnique(true).SetName(string(urlShortIndexName)),
		}

		_, err := urlCollection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("URL Short Index created successfully.")
	} else {
		fmt.Println("URL Short Index already exists.")
	}

	CurrentDb = &DB{User: userCollection, Url: urlCollection, Tracker: trackerCollection, Config: configCollection}
}
