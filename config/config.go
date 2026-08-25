// config/config.go
package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

func LoadConfig() (*Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	if _, err := time.LoadLocation(config.Database.Timezone); err != nil {
		return nil, fmt.Errorf("invalid DB_TIMEZONE %q: %w", config.Database.Timezone, err)
	}

	if strings.EqualFold(os.Getenv("GO_ENV"), "test") {
		applyTestDatabaseName(&config)
	}

	return &config, nil
}
