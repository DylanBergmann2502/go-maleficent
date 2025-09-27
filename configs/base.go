// configs/base.go
package configs

type ServerConfig struct {
	Host string `env:"SERVER_HOST,default=0.0.0.0"`
	Port int    `env:"SERVER_PORT,default=8000"`
}

type Config struct {
	Server ServerConfig
}

func GetBaseConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8000,
		},
	}
}
