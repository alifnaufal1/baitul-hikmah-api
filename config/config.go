package config

import (
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func LoadConfig() (*DatabaseConfig, error) {
	err := godotenv.Load()
  if err != nil {
    return nil, err
  }

  dbConfig := &DatabaseConfig{
    Host: getEnv("DB_HOST", "localhost"),
    Port: getEnv("DB_PORT", "5432"),
    User: getEnv("DB_USER", "postgres"),
    Password: getEnv("DB_PASSWORD", ""),
    DBName: getEnv("DB_NAME", "blog_db"),
  }

  return dbConfig, nil
}

func getEnv(key, defaultValue string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  }
  return defaultValue
}