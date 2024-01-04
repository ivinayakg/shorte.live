package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Name      string             `json:"name" validate:"required"`
	Email     string             `json:"email" validate:"required"`
	Picture   string             `json:"picture" validate:"required"`
	Token     string             `json:"token" bson:"-"`
	ID        primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"created_at" validate:"required"`
}

type URL struct {
	User        *User              `json:"user"`
	Destination string             `json:"destination"`
	Expiry      time.Time          `json:"expiry"`
	Short       string             `json:"short"`
	ID          primitive.ObjectID `json:"_id"`
	LastVisited time.Time          `json:"lastVisited"`
}

type URLDoc struct {
	User        primitive.ObjectID `json:"user,omitempty" bson:"user,omitempty"`
	Destination string             `json:"destination" validate:"required"`
	Expiry      time.Time          `json:"expiry" validate:"required"`
	Short       string             `json:"short" validate:"required"`
	ID          primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	LastVisited time.Time          `json:"lastVisited"`
}
