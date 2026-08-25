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
	if err := validateStorageCredentials(config.Storage); err != nil {
		return nil, err
	}
	if err := validateMailConfig(config.Mail); err != nil {
		return nil, err
	}

	if strings.EqualFold(os.Getenv("GO_ENV"), "test") {
		applyTestDatabaseName(&config)
	}

	return &config, nil
}

func validateMailConfig(mail MailConfig) error {
	if !mail.Configured() {
		return nil
	}
	if mail.Host == "" || mail.Port == 0 || mail.From == "" {
		return fmt.Errorf("MAIL_HOST, MAIL_PORT, and MAIL_FROM must be provided together")
	}
	if (mail.Username == "") != (mail.Password == "") {
		return fmt.Errorf("MAIL_USERNAME and MAIL_PASSWORD must be provided together")
	}

	return nil
}

func validateStorageCredentials(storage StorageConfig) error {
	if (storage.AccessKeyID == "") != (storage.SecretAccessKey == "") {
		return fmt.Errorf("S3_ACCESS_KEY_ID and S3_SECRET_ACCESS_KEY must be provided together")
	}

	return nil
}
