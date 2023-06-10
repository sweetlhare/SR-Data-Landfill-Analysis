package pgclient

import (
	"log"
	"svalka-service/internal/config"
)

var pgConfig config.PGConfig

// getPGConfig ...
func getPGConfig() config.PGConfig {
	if pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		pgConfig = cfg
	}

	return pgConfig
}
