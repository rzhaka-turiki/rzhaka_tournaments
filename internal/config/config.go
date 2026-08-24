package config

import (
	"errors"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server       ServerConfig
	Database     DatabaseConfig
	Results      ResultsConfig
	ApexVerifier ApexVerifierConfig
	Redis        RedisConfig
	Bot          BotConfig
	JWT          JWTConfig
}

type ApexVerifierConfig struct {
	GRPCAddr string
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	URL string
}

type ResultsConfig struct {
	GRPCAddr string
}

type RedisConfig struct {
	Addr string
}

type BotConfig struct {
	Token string
}

type JWTConfig struct {
	PublicKeyPath string
	AccessTTL     time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},

		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},

		Results: ResultsConfig{
			GRPCAddr: getEnv("GRPC_RESULTS_ADDR", "results:50051"),
		},

		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", ""),
		},

		Bot: BotConfig{
			Token: getEnv("BOT_TOKEN", ""),
		},

		JWT: JWTConfig{
			PublicKeyPath: getEnv("JWT_PUBLIC_KEY", "./keys/public.pem"),
			AccessTTL:     15 * time.Minute,
		},
		ApexVerifier: ApexVerifierConfig{
			GRPCAddr: getEnv(
				"GRPC_APEX_VERIFIER_ADDR",
				"localhost:50051",
			),
		},
	}
	if cfg.Database.URL == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
