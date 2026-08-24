// config/test.go
package config

import "strings"

func applyTestDatabaseName(config *Config) {
	if !strings.HasSuffix(config.Database.Database, "_test") {
		config.Database.Database += "_test"
	}
}
