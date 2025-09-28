// configs/local.go
package configs

func GetLocalConfig() *Config {
	config := GetBaseConfig()

	// Local environment overrides
	config.Log.Level = "debug"
	config.Log.Format = "console"

	return config
}
