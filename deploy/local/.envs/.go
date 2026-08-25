# deploy/local/.envs/.go
GO_ENV=local
# SERVER_HOST=0.0.0.0
# SERVER_PORT=8000
LOG_LEVEL=debug
LOG_FORMAT=console
# LOG_OUTPUT=stdout
PPROF_ENDPOINTS_ENABLED=true
# BACKGROUND_JOBS_CONCURRENCY=10

# CORS Configuration
# ------------------------------------------------------------------------------
# CORS_ALLOW_ORIGINS=*
# CORS_ALLOW_CREDENTIALS=false
# CORS_ALLOW_HEADERS=*

# Redis Configuration
# ------------------------------------------------------------------------------
REDIS_HOST=redis
REDIS_PORT=6379
# REDIS_PASSWORD=
# REDIS_DB=0

# Database Configuration
# ------------------------------------------------------------------------------
DB_HOST=postgres
DB_PORT=5432
DB_NAME=go_maleficent
DB_USER=fqrpqffwdjqsjhGgYdCShzeWJwGxDCyu
DB_PASSWORD=AuLjZ5WZzj4F5V7FAv6aRmkvuc3W8jhSq1ek7RWOfvWIgWtJ8tHVALmo9tVibJzJ
# DB_MAX_OPEN_CONNS=25
# DB_MAX_IDLE_CONNS=10
# DB_CONN_MAX_LIFETIME=4m
# DB_CONN_MAX_IDLE_TIME=15m
# DB_TIMEZONE=UTC
# POSTGRES_SSLMODE=disable

# S3 Storage (Garage)
# ------------------------------------------------------------------------------
S3_ENDPOINT=http://garage:3900
S3_REGION=garage
S3_BUCKET=local-go-maleficent
# S3_LOCATION=storage
S3_ACCESS_KEY_ID=GK5461d36b1ebf0cf4ee601aee
S3_SECRET_ACCESS_KEY=21022d6eec2f0020fee2c6f79544cbd620620b8542a71996a06d6bad9a5a6e0a
S3_USE_PATH_STYLE=true
