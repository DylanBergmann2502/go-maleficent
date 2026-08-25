# deploy/production/.envs/.go
GO_ENV=production
# SERVER_HOST=0.0.0.0
# SERVER_PORT=8000
LOG_LEVEL=info
LOG_FORMAT=json
# LOG_OUTPUT=stdout
PPROF_ENDPOINTS_ENABLED=false
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
DB_PASSWORD=bp8RBCtN99P62As7u4zgCdu7xc7bd9uCdT88KKWWmcTiDo0e9pd31JM6f51ntJrG
# DB_MAX_OPEN_CONNS=25
# DB_MAX_IDLE_CONNS=10
# DB_CONN_MAX_LIFETIME=4m
# DB_CONN_MAX_IDLE_TIME=15m
# DB_TIMEZONE=UTC
# POSTGRES_SSLMODE=disable

# S3 Storage
# ------------------------------------------------------------------------------
S3_ENDPOINT=
# S3_REGION=us-east-1
S3_BUCKET=production-go-maleficent
S3_ACCESS_KEY_ID=
S3_SECRET_ACCESS_KEY=
# S3_USE_PATH_STYLE=false
