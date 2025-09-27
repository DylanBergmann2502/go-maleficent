// configs/production.go
package configs

func GetProductionConfig() *Config {
	config := GetBaseConfig()

	// Production environment overrides would go here
	// For now, just return base config

	return config
}
