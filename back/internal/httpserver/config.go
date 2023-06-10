package httpserver

import (
	"log"
	"svalka-service/internal/config"
)

var (
	c config.HttpServerConfig
)

// getHttpServerConfig ...
func getHttpServerConfig() config.HttpServerConfig {
	if c == nil {
		cfg, err := config.NewHttpServerConfig()
		if err != nil {
			log.Fatal(ServerConfigError.AddDescription(err.Error()))
		}

		c = cfg
	}

	return c
}
