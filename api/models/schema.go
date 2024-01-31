package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnixTime int64

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
	LastVisited UnixTime           `json:"last_visited"`
	CreatedAt   UnixTime           `json:"created_at"`
}
