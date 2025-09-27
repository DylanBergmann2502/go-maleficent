// configs/local.go
package configs

func GetLocalConfig() *Config {
	config := GetBaseConfig()

	// Local environment overrides would go here
	// For now, just return base config

	return config
}
