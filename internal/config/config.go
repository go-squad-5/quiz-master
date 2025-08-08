package config

import (
	"fmt"

	"github.com/go-squad-5/quiz-master/internal/models"
)

type Config struct {
	Port string
	DSN  string
	// and more fields...
}

func Load() *Config {
	// TODO: Load configuration from environment variables, files, etc.
	// For simplicity, returning a hardcoded config.

	dbconfig := models.GetDBConfig()

	return &Config{
		Port: ":8080",
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dbconfig.User, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.DBName,
		),
	}
}
