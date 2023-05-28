package sessionrep

import (
	"log"
	"svalka-service/internal/config"
)

var redisConfig config.RedisRepConfig

// getRedisConfig ...
func getRedisConfig() config.RedisRepConfig {
	if redisConfig == nil {
		cfg, err := config.NewRedisRepConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		redisConfig = cfg
	}

	return redisConfig
}
