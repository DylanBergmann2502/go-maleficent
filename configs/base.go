// configs/base.go
package configs

type ServerConfig struct {
	Host string `env:"SERVER_HOST,default=0.0.0.0"`
	Port int    `env:"SERVER_PORT,default=8000"`
}

type LogConfig struct {
	Level      string   `env:"LOG_LEVEL,default=info"`
	Format     string   `env:"LOG_FORMAT,default=json"`
	OutputPath []string `env:"LOG_OUTPUT,default=stdout"`
}

type CORSConfig struct {
	AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS,default=*"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS,default=false"`
	AllowHeaders     []string `env:"CORS_ALLOW_HEADERS,default=*"`
}

type Config struct {
	Server ServerConfig
	Log    LogConfig
	CORS   CORSConfig
}

func GetBaseConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8000,
		},
		Log: LogConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: []string{"stdout"},
		},
		CORS: CORSConfig{
			AllowOrigins:     []string{"*"},
			AllowCredentials: false,
			AllowHeaders:     []string{"*"},
		},
	}
}
