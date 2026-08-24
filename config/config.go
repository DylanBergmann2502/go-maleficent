// config/config.go
package config

import (
	"os"
	"strings"

	"github.com/caarlos0/env/v11"
)

func LoadConfig() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}

	if strings.EqualFold(os.Getenv("GO_ENV"), "test") {
		applyTestDatabaseName(&config)
	}

	return &config, nil
}
