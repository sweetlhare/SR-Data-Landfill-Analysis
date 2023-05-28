package config

import (
	"os"

	"github.com/pkg/errors"
)

var _ AiClientConfig = (*aiClientConfig)(nil)

const (
	aiAddrEnv = "AI_ADDR"
)

type AiClientConfig interface {
	AiAddr() string
}

type aiClientConfig struct {
	aiAddr string
}

func NewAiClientConfig() (*aiClientConfig, error) {
	aiAddr := os.Getenv(aiAddrEnv)
	if aiAddr == "" {
		return nil, errors.New("ai client: addr not found")
	}

	return &aiClientConfig{
		aiAddr: aiAddr,
	}, nil
}

// AiAddr ...
func (cfg *aiClientConfig) AiAddr() string {
	return cfg.aiAddr
}
