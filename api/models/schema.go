package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnixTime int64
type CountryName string

type User struct {
	Name      string             `json:"name" validate:"required"`
	Email     string             `json:"email" validate:"required"`
	Picture   string             `json:"picture" validate:"required"`
	Token     string             `json:"token" bson:"-"`
	ID        primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	CreatedAt UnixTime           `json:"created_at" validate:"required"`
}

type URL struct {
	User        primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	UserDoc     *User              `json:"user_obj" bson:"-"`
	Destination string             `json:"destination" validate:"required"`
	Expiry      UnixTime           `json:"expiry" validate:"required"`
	Short       string             `json:"short" validate:"required"`
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UpdateAt    UnixTime           `json:"update_at" bson:"update_at"`
	CreatedAt   UnixTime           `json:"created_at" bson:"created_at"`
	TotalClicks int64              `json:"total_clicks" bson:"total_clicks"`
}

type RedirectEvent struct {
	URLId     primitive.ObjectID `json:"url_id,omitempty" bson:"url_id,omitempty"`
	URL       *URL               `json:"url_obj" bson:"-"`
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Geo       CountryName        `json:"geo" bson:"geo"`
	Device    string             `json:"device" bson:"device"`
	OS        string             `json:"os" bson:"os"`
	Referrer  string             `json:"referrer" bson:"referrer"`
	Timestamp UnixTime           `json:"timestamp" bson:"timestamp"`
}
