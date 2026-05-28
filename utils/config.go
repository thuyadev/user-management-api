package utils

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppPort string
	AppEnv  string

	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	MongoURI string
	MongoDB  string

	JWTSecret      string
	JWTExpiryHours int

	AdminName     string
	AdminEmail    string
	AdminPassword string

	AIEnabled bool
	AIAPIKey  string
	AIAPIURL  string
	AIModel   string

	APIKeyRequired     bool
	APIKey             string
	APIKeyHeader       string
	CORSAllowedOrigins []string

	SwaggerEnabled  bool
	SwaggerPath     string
	SwaggerUser     string
	SwaggerPassword string
}

func LoadConfig() *Config {
	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "user_management"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "user_management_logs"),

		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 24),

		AdminName:     getEnv("ADMIN_NAME", "Admin User"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "password123"),

		AIEnabled: getEnvBool("AI_ENABLED", false),
		AIAPIKey:  getEnv("AI_API_KEY", ""),
		AIAPIURL:  getEnv("AI_API_URL", "https://api.openai.com/v1/chat/completions"),
		AIModel:   getEnv("AI_MODEL", "gpt-4o-mini"),

		APIKeyRequired:     getEnvBool("API_KEY_REQUIRED", true),
		APIKey:             getEnv("API_KEY", ""),
		APIKeyHeader:       getEnv("API_KEY_HEADER", "X-API-Key"),
		CORSAllowedOrigins: splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "")),

		SwaggerEnabled:  getEnvBool("SWAGGER_ENABLED", true),
		SwaggerPath:     getEnv("SWAGGER_PATH", "/swagger"),
		SwaggerUser:     getEnv("SWAGGER_USER", "swagger"),
		SwaggerPassword: getEnv("SWAGGER_PASSWORD", "swagger"),
	}
}

func splitCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			out = append(out, s)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}
