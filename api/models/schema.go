package models

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	Name      string `bun:"name, notnull"`
	Email     string `bun:"email, notnull"`
	Picture   string `bun:"picture, notnull"`
	Token     string
	ID        string    `bun:"id,pk,type:uuid, default:gen_random_uuid()"`
	CreatedAt time.Time `bun:"created_at, notnull"`
}
