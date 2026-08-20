// configs/config.go
package configs

import (
	"context"
	"os"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

func LoadConfig() (*Config, error) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	var config *Config

	switch env {
	case "local":
		config = GetLocalConfig()
	case "production":
		config = GetProductionConfig()
	case "test":
		config = GetTestConfig()
	default:
		config = GetLocalConfig()
	}

	// sethvargo/go-envconfig only overwrites fields where env vars exist
	// This preserves your local overrides unless explicitly overridden
	err := envconfig.Process(context.Background(), config)
	if err != nil {
		return nil, err
	}

	if env == "test" && !strings.HasSuffix(config.Database.Database, "_test") {
		config.Database.Database += "_test"
	}

	return config, nil
}
