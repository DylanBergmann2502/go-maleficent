// configs/test.go
package configs

func GetTestConfig() *Config {
	config := GetBaseConfig()

	// Test environment overrides would go here
	// For now, just return base config

	return config
}
