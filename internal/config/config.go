package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL   string
	EncryptionKey string // base64

	// Auth
	AuthPort            string
	DiscordClientID     string
	DiscordClientSecret string
	AuthRedirectURL     string
	FrontendURL         string
	JWTSecret           string
	JWTAccessTTL        time.Duration
	JWTRefreshTTL       time.Duration

	// main api links
	ServerPort      string
	GRPCResultsAddr string
	RedisAddr       string
	BotToken        string // JWT-token for discord bot
}

func Load() *Config {
	// env variables
	return &Config{
		DatabaseURL:         getEnv("DATABASE_URL", "postgres://user:pass@localhost/apex?sslmode=disable"),
		EncryptionKey:       getEnv("ENCRYPTION_KEY", ""),
		AuthPort:            getEnv("AUTH_PORT", "8081"),
		DiscordClientID:     getEnv("DISCORD_CLIENT_ID", ""),
		DiscordClientSecret: getEnv("DISCORD_CLIENT_SECRET", ""),
		AuthRedirectURL:     getEnv("AUTH_REDIRECT_URL", "http://localhost:8081/auth/discord/callback"),
		FrontendURL:         getEnv("FRONTEND_URL", "http://localhost:3000"),
		JWTSecret:           getEnv("JWT_SECRET", "change-me"),
		JWTAccessTTL:        15 * time.Minute,
		JWTRefreshTTL:       7 * 24 * time.Hour,
		ServerPort:          getEnv("SERVER_PORT", "8080"),
		GRPCResultsAddr:     getEnv("GRPC_RESULTS_ADDR", "results:50051"),
		RedisAddr:           getEnv("REDIS_ADDR", ""),
		BotToken:            getEnv("BOT_TOKEN", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
