package config

import (
	"os"
	"time"

	"github.com/pkg/errors"
)

var _ HttpServerConfig = (*httpServerConfig)(nil)

const (
	httpServerHostEnvName            = "HTTPSERVER_HOST"
	httpServerModeEnvName            = "HTTPSERVER_MODE"
	httpServerShutdownTimeoutEnvName = "HTTPSERVER_SHUTDOWNTIMEOUT"
)

type HttpServerConfig interface {
	Host() string
	Mode() string
	ShutdownTimeout() time.Duration
}

type httpServerConfig struct {
	host            string
	mode            string
	shutdownTimeout time.Duration
}

func NewHttpServerConfig() (*httpServerConfig, error) {
	host := os.Getenv(httpServerHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http server: host not found")
	}

	mode := os.Getenv(httpServerModeEnvName)
	if len(host) == 0 {
		return nil, errors.New("http server: mode not found")
	}

	shutdownTimeout, err := time.ParseDuration(os.Getenv(httpServerHostEnvName))
	if len(host) == 0 {
		return nil, errors.Errorf("http server: shutdown timeout not found, %s", err.Error())
	}

	return &httpServerConfig{
		host:            host,
		mode:            mode,
		shutdownTimeout: shutdownTimeout,
	}, nil
}

func (cfg *httpServerConfig) Host() string {
	return cfg.host
}

func (cfg *httpServerConfig) Mode() string {
	return cfg.mode
}

func (cfg *httpServerConfig) ShutdownTimeout() time.Duration {
	return cfg.shutdownTimeout
}
