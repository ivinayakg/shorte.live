package helpers

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisDB struct {
	Client *redis.Client
}

var Redis redisDB

func RedisSetup() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	Redis.Client = redis.NewClient(opt)
}

func (r redisDB) SetJSON(key string, value interface{}, expiry time.Duration) error {
	p, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(context.Background(), key, p, expiry).Err()
}

func (r redisDB) GetJSON(key string, dest interface{}) error {
	p, err := r.Client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(p), dest)
}
