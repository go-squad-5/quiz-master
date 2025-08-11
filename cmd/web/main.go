package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-squad-5/quiz-master/internal/app"
	"github.com/go-squad-5/quiz-master/internal/config"
	"github.com/go-squad-5/quiz-master/internal/repositories"
	"github.com/go-squad-5/quiz-master/internal/router"
)

func main() {
	cfg := config.Load()
	// Initialize the application
	app := app.NewApp(cfg)
	// Initialize the router with the app instance
	// Start the HTTP server
	app.Repository = InitDB(app) // InitDB(app)
	log.Println("Connected to MySQL!")

	router := router.NewRouter(app)
	route := router.Route()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(app.Config.Port, route))
}

func InitDB(app *app.App) *repositories.Repository {
	respository, err := repositories.NewMySQLRepository(app.Config.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	return &respository
}
