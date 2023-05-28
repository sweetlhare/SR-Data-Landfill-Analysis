package config

import (
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var _ RedisRepConfig = (*redisRepConfig)(nil)

const (
	redisAddrEnv           = "REDIS_ADDR"
	redisPasswordEnv       = "REDIS_PASSWORD"
	redisDbEnv             = "REDIS_DB"
	redisSessionTimeoutEnv = "REDIS_SESSION_TIMEOUT"
)

type RedisRepConfig interface {
	RedisAddr() string
	RedisPassword() string
	RedisDb() int64
	RedisSessionTimeout() time.Duration
}

type redisRepConfig struct {
	redisAddr           string
	redisPassword       string
	redisDb             int64
	redisSessionTimeout time.Duration
}

func NewRedisRepConfig() (*redisRepConfig, error) {
	redisAddr := os.Getenv(redisAddrEnv)
	if redisAddr == "" {
		return nil, errors.New("redis rep: addr not found")
	}

	redisPassword := os.Getenv(redisPasswordEnv)

	redisDbString := os.Getenv(redisDbEnv)
	redisDb, err := strconv.ParseInt(redisDbString, 10, 64)
	if err != nil {
		return nil, errors.Errorf("redis rep: redis db not found, %s", err.Error())
	}

	redisSessionTimeout, err := time.ParseDuration(os.Getenv(redisSessionTimeoutEnv))
	if err != nil {
		return nil, errors.Errorf("redis rep: shutdown timeout not found, %s", err.Error())
	}

	return &redisRepConfig{
		redisAddr:           redisAddr,
		redisPassword:       redisPassword,
		redisDb:             redisDb,
		redisSessionTimeout: redisSessionTimeout,
	}, nil
}

// RedisAddr ...
func (cfg *redisRepConfig) RedisAddr() string {
	return cfg.redisAddr
}

// RedisPassword ...
func (cfg *redisRepConfig) RedisPassword() string {
	return cfg.redisPassword
}

// RedisDb ...
func (cfg *redisRepConfig) RedisDb() int64 {
	return cfg.redisDb
}

// RedisSessionTimeout ...
func (cfg *redisRepConfig) RedisSessionTimeout() time.Duration {
	return cfg.redisSessionTimeout
}
