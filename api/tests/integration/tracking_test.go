package integration_tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ivinayakg/shorte.live/api/helpers"
	"github.com/ivinayakg/shorte.live/api/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestURLRedirectTracking(t *testing.T) {
	resp, err := RedirecthttpClient.Get(ServerURL + "/" + URLFixture.Short)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	destinationURL := URLFixture.Destination

	assert.Equal(t, resp.StatusCode, http.StatusMovedPermanently, "Excpected status code to be 301")
	assert.Contains(t, resp.Header.Get("Location"), destinationURL, "Expected redirect to destination url")

	var result *models.RedirectEvent = nil

	for result == nil {
		err := helpers.CurrentDb.RedirectEvent.FindOne(context.Background(), bson.M{"url_id": URLFixture.ID}).Decode(&result)
		if err != nil && err != mongo.ErrNoDocuments {
			t.Log(err)
			t.Fail()
		}
		time.Sleep(time.Second * 2)
	}

	assert.Equal(t, (*result).URLId, URLFixture.ID, "Expected URL ID to be the same")
}
