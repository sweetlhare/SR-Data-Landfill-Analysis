package aiclient

import (
	"log"
	"svalka-service/internal/config"
)

var AiClientConfig config.AiClientConfig

// getAiClientConfig ...
func getAiClientConfig() config.AiClientConfig {
	if AiClientConfig == nil {
		cfg, err := config.NewAiClientConfig()
		if err != nil {
			log.Fatalf("failed to get ai client config: %s", err.Error())
		}

		AiClientConfig = cfg
	}

	return AiClientConfig
}
