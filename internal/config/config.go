package config

import "fmt"

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		User:     "quizuser",
		Password: "quizpass",
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

	return &Config{
		Port: ":8090",
		DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dbconfig.User, dbconfig.Password, dbconfig.Host, dbconfig.Port, dbconfig.DBName,
		),
	}
}
