// configs/test.go
package configs

import "strings"

func GetTestConfig() *Config {
	config := GetBaseConfig()

	return config
}

func applyTestDatabaseName(config *Config) {
	if !strings.HasSuffix(config.Database.Database, "_test") {
		config.Database.Database += "_test"
	}
}
