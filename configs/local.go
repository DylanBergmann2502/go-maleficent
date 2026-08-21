// configs/local.go
package configs

func GetLocalConfig() *Config {
	config := GetBaseConfig()

	// Local environment overrides
	config.Log.Level = "debug"
	config.Log.Format = "console"
	config.Debug.PprofEndpointsEnabled = true

	return config
}
