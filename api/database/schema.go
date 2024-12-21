package database

import (
	"database/sql/driver"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnixTime int64
type CountryName string

func (u UnixTime) Value() (driver.Value, error) {
	return time.Unix(int64(u), 0), nil
}

func (ut *UnixTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case time.Time:
		*ut = UnixTime(v.Unix())
		return nil
	default:
		return fmt.Errorf("cannot convert %v to UnixTime", value)
	}
}

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
	Temporary   bool               `json:"temporary" bson:"temporary"`
}

type ClickEvent struct {
	URLId     string      `json:"url_id,omitempty"`
	Geo       CountryName `json:"geo"`
	Device    string      `json:"device"`
	OS        string      `json:"os"`
	Referrer  string      `json:"referrer"`
	Timestamp UnixTime    `json:"timestamp"`
}
