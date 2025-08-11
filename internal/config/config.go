package config

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		User:     "root",
		Password: "your_secure_password",
		Host:     "localhost",
		Port:     "3306",
		DBName:   "quizdb",
	}
}

type Config struct {
	Port string
	DSN  string
	// and more fields...
}

func Load() *Config {
	// TODO: Load configuration from environment variables, files, etc.
	// For simplicity, returning a hardcoded config.

	dbconfig := GetDBConfig()

	dsn := os.Getenv("DB_DSN")

	if dsn != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dbconfig.User, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.DBName,
		)
	}

	return &Config{
		Port: ":8090",
		DSN:  dsn,
	}
}
