package testhelper

import (
	"time"

	"github.com/ivinayakg/shorte.live/api/helpers"
)

func PutSystemUnderMaintenance(redis *helpers.RedisDB, val bool) {
	config := helpers.GetDefaultSystemConfig()
	config.Maintenance = val
	redis.SetJSON(string(helpers.SystemConfigNameCacheKey), config, time.Hour*24)
}
