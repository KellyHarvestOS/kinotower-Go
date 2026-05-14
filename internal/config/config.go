package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppAddr         string
	Env             string
	DatabaseURL     string
	JWTSecret       string
	JWTTTL          time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
	RateLimitRPM    int
	AutoMigrate     bool
	AutoSeed        bool
	MigrateOnly     bool
	SeedOnly        bool
}

func Load() Config {
	_ = godotenv.Load()

	dbURL := env("DATABASE_URL", "")
	if dbURL == "" {
		dbURL = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=%s",
			env("POSTGRES_USER", "postgres"),
			env("POSTGRES_PASSWORD", ""),
			env("POSTGRES_HOST", "localhost"),
			env("POSTGRES_PORT", "5432"),
			env("POSTGRES_DB", "kinotower"),
			env("DB_SSLMODE", "disable"),
		)
	}

	return Config{
		AppAddr:         env("APP_ADDR", ":8080"),
		Env:             env("APP_ENV", "local"),
		DatabaseURL:     dbURL,
		JWTSecret:       env("JWT_SECRET", "change-me-in-env"),
		JWTTTL:          durationEnv("JWT_TTL", 24*time.Hour),
		ShutdownTimeout: durationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		RequestTimeout:  durationEnv("REQUEST_TIMEOUT", 15*time.Second),
		RateLimitRPM:    intEnv("RATE_LIMIT_RPM", 120),
		AutoMigrate:     boolEnv("AUTO_MIGRATE", true),
		AutoSeed:        boolEnv("AUTO_SEED", false),
		MigrateOnly:     boolEnv("MIGRATE_ONLY", false),
		SeedOnly:        boolEnv("SEED_ONLY", false),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func intEnv(key string, fallback int) int {
	value, err := strconv.Atoi(env(key, ""))
	if err != nil {
		return fallback
	}
	return value
}

func boolEnv(key string, fallback bool) bool {
	value, err := strconv.ParseBool(env(key, ""))
	if err != nil {
		return fallback
	}
	return value
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	value, err := time.ParseDuration(env(key, ""))
	if err != nil {
		return fallback
	}
	return value
}
