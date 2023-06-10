package config

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
)

func Init(_ context.Context) error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("init config error: %s", err.Error())
	}
	return nil
}
