package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
	Upload   UploadConfig
}

// ServerConfig 服务配置
type ServerConfig struct {
	Host string
	Port string
	Mode string
}

// PostgresConfig PostgreSQL配置
type PostgresConfig struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
	Timezone string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string
	ExpireHours     int
	RefreshExpireHours int
}

// LogConfig 日志配置
type LogConfig struct {
	Level    string
	FilePath string
}

// UploadConfig 文件上传配置
type UploadConfig struct {
	MaxSize     int64
	AllowedTypes string
	SavePath    string
}

// Load 加载配置
func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Postgres: PostgresConfig{
			Host:     getEnv("PG_HOST", "127.0.0.1"),
			Port:     getEnv("PG_PORT", "5432"),
			DBName:   getEnv("PG_DBNAME", "genealogy"),
			User:     getEnv("PG_USER", "postgres"),
			Password: getEnv("PG_PASSWORD", ""),
			SSLMode:  getEnv("PG_SSL_MODE", "disable"),
			Timezone: getEnv("PG_TIMEZONE", "Asia/Shanghai"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", "default-secret"),
			ExpireHours:         getEnvInt("JWT_EXPIRE_HOURS", 24),
			RefreshExpireHours:  getEnvInt("JWT_REFRESH_EXPIRE_HOURS", 168),
		},
		Log: LogConfig{
			Level:    getEnv("LOG_LEVEL", "info"),
			FilePath: getEnv("LOG_FILE_PATH", "./logs/app.log"),
		},
		Upload: UploadConfig{
			MaxSize:      int64(getEnvInt("UPLOAD_MAX_SIZE", 10485760)),
			AllowedTypes: getEnv("UPLOAD_ALLOWED_TYPES", "jpg,jpeg,png,gif,pdf"),
			SavePath:     getEnv("UPLOAD_SAVE_PATH", "./uploads"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// DSN 返回PostgreSQL连接字符串
func (p *PostgresConfig) DSN() string {
	return "host=" + p.Host +
		" user=" + p.User +
		" password=" + p.Password +
		" dbname=" + p.DBName +
		" port=" + p.Port +
		" sslmode=" + p.SSLMode +
		" TimeZone=" + p.Timezone
}

// Addr 返回Redis地址
func (r *RedisConfig) Addr() string {
	return r.Host + ":" + r.Port
}
