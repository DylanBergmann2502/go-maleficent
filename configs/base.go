// configs/base.go
package configs

import "time"

type ServerConfig struct {
	Host string `env:"SERVER_HOST,overwrite,default=0.0.0.0"`
	Port int    `env:"SERVER_PORT,overwrite,default=8000"`
}

type DatabaseConfig struct {
	// Basic connection
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Database string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`

	// Essential pool settings (prevent PostgreSQL exhaustion)
	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS,overwrite,default=25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS,overwrite,default=10"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME,overwrite,default=4m"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME,overwrite,default=15m"`

	// Explicit settings
	Timezone string `env:"DB_TIMEZONE,overwrite,default=UTC"`
	SSLMode  string `env:"POSTGRES_SSLMODE,overwrite,default=disable"`
}

type LogConfig struct {
	Level      string   `env:"LOG_LEVEL,overwrite,default=info"`
	Format     string   `env:"LOG_FORMAT,overwrite,default=json"`
	OutputPath []string `env:"LOG_OUTPUT,overwrite,default=stdout"`
}

type CORSConfig struct {
	AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS,overwrite,default=*"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS,overwrite,default=false"`
	AllowHeaders     []string `env:"CORS_ALLOW_HEADERS,overwrite,default=*"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
	CORS     CORSConfig
}

func GetBaseConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8000,
		},
		Database: DatabaseConfig{
			MaxOpenConns:    25,
			MaxIdleConns:    10,
			ConnMaxLifetime: 4 * time.Minute,
			ConnMaxIdleTime: 15 * time.Minute,
			Timezone:        "UTC",
			SSLMode:         "disable",
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
