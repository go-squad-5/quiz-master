package models

// package config

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
