package config

import "os"

type Config struct {
	DBDSN string
	Port  string
}

func LoadConfig() *Config {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = os.Getenv("DATABASE_URL")
	}
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5432/finops?sslmode=disable"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DBDSN: dsn,
		Port:  port,
	}
}
