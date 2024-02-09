package helpers

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TrackEventType string

const track_event_redis_key = "track_event"

type TrackerType struct {
	maxEvents      int
	eventsLength   int
	flushFrequency time.Duration
	mutex          sync.Mutex
}

var Tracker *TrackerType

func flushToDB(events []*bson.M) {
	var redirectEvents []interface{}
	for _, data := range events {
		urlOID, _ := (*data)["url_id"].(primitive.ObjectID)
		(*data)["url_id"] = urlOID
		redirectEvents = append(redirectEvents, data)
	}

	if len(redirectEvents) > 0 {
		_, err := CurrentDb.RedirectEvent.InsertMany(context.TODO(), redirectEvents)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (eq *TrackerType) flush() {
	eq.mutex.Lock()
	defer eq.mutex.Unlock()

	jsondata, err := Redis.Client.LRange(context.TODO(), track_event_redis_key, 0, int64(eq.eventsLength)).Result()
	if err != nil {
		log.Fatal(err)
	}

	Redis.Client.LTrim(context.TODO(), track_event_redis_key, int64(eq.eventsLength), -1)
	eq.eventsLength = 0

	var data []*bson.M

	for _, v := range jsondata {
		var temp bson.M
		bson.Unmarshal([]byte(v), &temp)
		data = append(data, &temp)
	}

	go flushToDB(data)
}

func (eq *TrackerType) CaptureRedirectEvent(device string, geo string, os string, referrer string, urlId primitive.ObjectID, timestamp int64) {
	// Your slice
	data := bson.M{
		"device":    device,
		"geo":       geo,
		"os":        os,
		"referrer":  referrer,
		"url_id":    urlId,
		"timestamp": timestamp,
	}

	jsonData, _ := bson.Marshal(data)

	// Push the entire slice as a single element into the Redis list
	err := Redis.Client.LPush(context.TODO(), track_event_redis_key, jsonData).Err()
	if err != nil {
		log.Fatal(err)
	}

	eq.eventsLength += 1
	if eq.eventsLength >= eq.maxEvents {
		eq.flush()
	}
}

func (eq *TrackerType) StartFlush() {
	ticker := time.NewTicker(eq.flushFrequency)
	defer ticker.Stop()

	for range ticker.C {
		eq.flush()
	}
}

func SetupTracker(dur time.Duration, maxEvents int, eventsLength int) {
	Redis.Client.LTrim(context.TODO(), track_event_redis_key, int64(eventsLength), -1)
	Tracker = &TrackerType{
		maxEvents:      maxEvents,
		eventsLength:   eventsLength,
		flushFrequency: dur,
	}
}
