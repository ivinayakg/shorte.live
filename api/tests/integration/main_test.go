package integration_tests

import (
	"context"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/shorte.live/api/controllers"
	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/middleware"
	"github.com/ivinayakg/shorte.live/api/models"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ServerURL    string
	TestDb       *helpers.DB
	TestRedis    *helpers.RedisDB
	TestDbClient *mongo.Client
)

// Test remaining for
// - Ratelimiting
// - System Maintenance

var UserFixture1 = models.User{Name: "Test User 1", Email: "test1@gmail.com", Picture: "https://lh3.googleusercontent.com/a-/AOh14Gh"}
var UserFixture2 = models.User{Name: "Test User 2", Email: "test2@gmail.com", Picture: "https://lh3.googleusercontent.com/a-/AOh14Gh"}

var URLFixture = &models.URL{User: primitive.NilObjectID, Destination: "https://www.google.com", Expiry: models.UnixTime(time.Now().Add(time.Hour * 5).Unix()), Short: "test", UpdateAt: models.UnixTime(time.Now().Unix()), CreatedAt: models.UnixTime(time.Now().Unix())}
var ExpiredURLFixture = &models.URL{User: primitive.NilObjectID, Destination: "https://www.google.com", Expiry: models.UnixTime(time.Now().Add(-time.Hour).Unix()), Short: "test_expired", UpdateAt: models.UnixTime(time.Now().Unix()), CreatedAt: models.UnixTime(time.Now().Unix())}

func TestMain(m *testing.M) {
	// Set up HTTP server
	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()
	ServerURL = server.URL
	err := godotenv.Load("../test.env")
	if err != nil {
		fmt.Println(err)
	}
	helpers.CreateDBInstance()
	helpers.RedisSetup()

	TestDb = helpers.CurrentDb
	TestRedis = helpers.Redis
	TestDbClient = helpers.DBClient
	// clean up database after tests
	defer TestDbClient.Database(os.Getenv("DB_NAME")).Drop(context.Background())
	defer TestRedis.Client.FlushAll(context.Background())

	CreateFixtures(TestDb)

	// Run integration tests
	exitCode := m.Run()
	os.Exit(exitCode)
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(middleware.Authentication)

	// user routes
	router.HandleFunc("/user/sign_in_with_google", controllers.SignInWithGoogle).Methods("GET")
	protectedRouter.HandleFunc("/user/self", controllers.SelfUser).Methods("GET")

	// url resolve routes
	router.HandleFunc("/{short}", controllers.ResolveURL).Methods("GET")

	// url routes
	protectedRouter.HandleFunc("/url", controllers.ShortenURL).Methods("POST")
	protectedRouter.HandleFunc("/url/all", controllers.GetUserURL).Methods("GET")
	protectedRouter.HandleFunc("/url/{id}", controllers.UpdateUrl).Methods("PATCH")
	protectedRouter.HandleFunc("/url/{id}", controllers.DeleteUrl).Methods("DELETE")

	// system routes
	router.HandleFunc("/system/available", controllers.SystemAvailable).Methods("GET")

	return router
}

func CreateFixtures(db *helpers.DB) {
	userRes, _ := db.User.InsertOne(context.Background(), UserFixture2)
	UserFixture2.ID = userRes.InsertedID.(primitive.ObjectID)
	userRes, _ = db.User.InsertOne(context.Background(), UserFixture1)
	UserFixture1.ID = userRes.InsertedID.(primitive.ObjectID)

	URLFixture.User = userRes.InsertedID.(primitive.ObjectID)
	ExpiredURLFixture.User = userRes.InsertedID.(primitive.ObjectID)

	URLFixture.User = UserFixture1.ID
	ExpiredURLFixture.User = UserFixture1.ID

	urlRes, _ := db.Url.InsertOne(context.Background(), URLFixture)
	URLFixture.ID = urlRes.InsertedID.(primitive.ObjectID)

	urlRes, _ = db.Url.InsertOne(context.Background(), ExpiredURLFixture)
	ExpiredURLFixture.ID = urlRes.InsertedID.(primitive.ObjectID)
}
