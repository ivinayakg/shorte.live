package helpers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLLimit struct {
	Value  int `json:"value"`
	Expiry int `json:"expiry"`
}

type RateConfig struct {
	Name  string               `json:"name"`
	Limit map[string]*URLLimit `json:"limit"`
	ID    primitive.ObjectID   `json:"_id,omitempty"  bson:"_id,omitempty"`
}

type RateLimitLog struct {
	Used     int           `json:"used"`
	Limit    int           `json:"limit"`
	CoolDown time.Duration `json:"cooldown"`
}

const RateConfigName string = "rate_limit_config"
const RateConfigNameCacheKey string = "cached_rate_limit_config"

func getDefaultRateConfig() RateConfig {
	limits := map[string]*URLLimit{"*": {Value: 10, Expiry: 30}}

	return RateConfig{
		Name:  RateConfigName,
		Limit: limits,
	}
}

func GetRateConfig(revalidateCache bool) *RateConfig {
	rateConfig := &RateConfig{}

	err := Redis.GetJSON(RateConfigNameCacheKey, rateConfig)
	if err != nil {
		fmt.Println(err)
	}

	if rateConfig.ID != primitive.NilObjectID && !revalidateCache {
		return rateConfig
	} else {
		var rateConfigFilters = bson.M{"name": RateConfigName}

		err = CurrentDb.Config.FindOne(context.Background(), rateConfigFilters).Decode(rateConfig)
		if err != nil && err != mongo.ErrNoDocuments {
			fmt.Println(err)
			return nil
		}

		if err == mongo.ErrNoDocuments {
			defaultRateConfig := getDefaultRateConfig()
			res, err := CurrentDb.Config.InsertOne(context.Background(), defaultRateConfig)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			rateConfig = &defaultRateConfig
			rateConfig.ID = res.InsertedID.(primitive.ObjectID)
		}

		Redis.SetJSON(RateConfigNameCacheKey, rateConfig, time.Duration(time.Minute*60))
		return rateConfig
	}
}

func RateLimit(r *http.Request, auth string, defaultLimit *URLLimit) (time.Duration, error) {
	rateConfig := GetRateConfig(false)
	urlRateName := r.URL.Path + "-" + r.Method

	urlRateConfig, found := rateConfig.Limit[urlRateName]
	if !found {
		urlRateConfig = rateConfig.Limit["*"]
	}

	if defaultLimit != nil {
		urlRateConfig = defaultLimit
	}

	if auth == "" {
		auth = GetUserIP(r)
	}

	rateLimitLog := &RateLimitLog{}
	userRateLimitKey := auth + "-" + urlRateName
	c := context.Background()

	cacheExpiry, err := Redis.Client.TTL(c, userRateLimitKey).Result()
	if err != nil {
		return 0, err
	}

	if cacheExpiry <= 0 {
		rateLimitLog = &RateLimitLog{
			CoolDown: time.Duration(urlRateConfig.Expiry) * time.Minute,
			Limit:    urlRateConfig.Value,
			Used:     1,
		}
	} else {
		err = Redis.GetJSON(userRateLimitKey, rateLimitLog)
		if err != nil {
			return 0, err
		}
		rateLimitLog.CoolDown = cacheExpiry

		if rateLimitLog.Limit <= rateLimitLog.Used {
			return rateLimitLog.CoolDown, fmt.Errorf("too many requests")
		}

		rateLimitLog.Used += 1
	}

	if err := Redis.SetJSON(userRateLimitKey, rateLimitLog, rateLimitLog.CoolDown); err != nil {
		return 0, err
	}

	return 0, nil
}
