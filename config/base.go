// config/base.go
package config

import "time"

type ServerConfig struct {
	Host string `env:"SERVER_HOST" envDefault:"0.0.0.0"`
	Port int    `env:"SERVER_PORT" envDefault:"8000"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Database string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`

	MaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" envDefault:"10"`
	ConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" envDefault:"4m"`
	ConnMaxIdleTime time.Duration `env:"DB_CONN_MAX_IDLE_TIME" envDefault:"15m"`
	Timezone        string        `env:"DB_TIMEZONE" envDefault:"UTC"`
	SSLMode         string        `env:"POSTGRES_SSLMODE" envDefault:"disable"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST,required"`
	Port     int    `env:"REDIS_PORT,required"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}

type BackgroundJobsConfig struct {
	Concurrency int `env:"BACKGROUND_JOBS_CONCURRENCY" envDefault:"10"`
}

type LogConfig struct {
	Level      string   `env:"LOG_LEVEL" envDefault:"info"`
	Format     string   `env:"LOG_FORMAT" envDefault:"json"`
	OutputPath []string `env:"LOG_OUTPUT" envDefault:"stdout"`
}

type CORSConfig struct {
	AllowOrigins     []string `env:"CORS_ALLOW_ORIGINS" envDefault:"*"`
	AllowCredentials bool     `env:"CORS_ALLOW_CREDENTIALS" envDefault:"false"`
	AllowHeaders     []string `env:"CORS_ALLOW_HEADERS" envDefault:"*"`
}

type DebugConfig struct {
	PprofEndpointsEnabled bool `env:"PPROF_ENDPOINTS_ENABLED" envDefault:"false"`
}

type StorageConfig struct {
	Endpoint        string `env:"S3_ENDPOINT"`
	Region          string `env:"S3_REGION" envDefault:"us-east-1"`
	Bucket          string `env:"S3_BUCKET,required"`
	Location        string `env:"S3_LOCATION" envDefault:"storage"`
	AccessKeyID     string `env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY"`
	UsePathStyle    bool   `env:"S3_USE_PATH_STYLE" envDefault:"false"`
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Jobs     BackgroundJobsConfig
	Log      LogConfig
	CORS     CORSConfig
	Debug    DebugConfig
	Storage  StorageConfig
}
