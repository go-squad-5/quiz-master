package app

import (
	"github.com/go-squad-5/quiz-master/internal/config"
	"github.com/go-squad-5/quiz-master/internal/repositories"
)

type App struct {
	Config     *config.Config
	Repository *repositories.Repository // and more fields... db, logger, etc.
}

func NewApp(config *config.Config) *App {
	return &App{
		Config: config,
		// Initialize other fields like db, logger, etc.
	}
}
