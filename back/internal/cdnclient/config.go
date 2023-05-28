package cdnclient

import (
	"log"
	"svalka-service/internal/config"
)

var ftpClientConfig config.FtpClientConfig

// getFtpClientConfig ...
func getFtpClientConfig() config.FtpClientConfig {
	if ftpClientConfig == nil {
		cfg, err := config.NewFtpClientConfig()
		if err != nil {
			log.Fatalf("failed to get ftp client config: %s", err.Error())
		}

		ftpClientConfig = cfg
	}

	return ftpClientConfig
}
