package helpers

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SystemConfigKey string

type SystemConfig struct {
	Name        SystemConfigKey    `json:"name"`
	Maintenance bool               `json:"maintenance"`
	ID          primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
}

const SystemConfigName SystemConfigKey = "system_limit_config"
const SystemConfigNameCacheKey SystemConfigKey = "cached_system_limit_config"

func GetDefaultSystemConfig() *SystemConfig {
	return &SystemConfig{
		Name:        SystemConfigName,
		Maintenance: false,
	}
}

func GetSystemConfig(revalidateCache bool) *SystemConfig {
	systemConfig := &SystemConfig{}

	err := Redis.GetJSON(string(SystemConfigNameCacheKey), systemConfig)
	if err != nil {
		fmt.Println(err)
	}

	if systemConfig.ID != primitive.NilObjectID && !revalidateCache {
		return systemConfig
	} else {
		var systemConfigFilters = bson.M{"name": SystemConfigName}

		err = CurrentDb.Config.FindOne(context.Background(), systemConfigFilters).Decode(systemConfig)
		if err != nil && err != mongo.ErrNoDocuments {
			fmt.Println(err)
			return nil
		}

		if err == mongo.ErrNoDocuments {
			defaultSystemConfig := GetDefaultSystemConfig()
			res, err := CurrentDb.Config.InsertOne(context.Background(), defaultSystemConfig)
			if err != nil {
				fmt.Println(err)
				return nil
			}

			systemConfig = defaultSystemConfig
			systemConfig.ID = res.InsertedID.(primitive.ObjectID)
		}

		Redis.SetJSON(string(SystemConfigNameCacheKey), systemConfig, time.Duration(time.Minute*60))
		return systemConfig
	}
}

func SystemUnderMaintenance(revalidate bool) bool {
	systemConfig := GetSystemConfig(revalidate)
	if systemConfig == nil {
		return false
	}

	return systemConfig.Maintenance
}
