package app

import (
	"database/sql"

	"github.com/go-squad-5/quiz-master/internal/config"
)

type App struct {
	Config *config.Config
	Db     *sql.DB
	// and more fields... db, logger, etc.
}

func NewApp(config *config.Config) *App {
	return &App{
		Config: config,
		// Initialize other fields like db, logger, etc.
	}
}
