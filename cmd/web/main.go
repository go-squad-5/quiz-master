package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-squad-5/quiz-master/internal/app"
	"github.com/go-squad-5/quiz-master/internal/config"
	"github.com/go-squad-5/quiz-master/internal/router"
)

func main() {
	cfg := config.Load()
	// Initialize the application
	app := app.NewApp(cfg)
	// Initialize the router with the app instance
	route := router.Route()
	// Start the HTTP server
	InitDB(app)
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(app.Config.Port, route))
}

func InitDB(app *app.App) {
	// cfg := config.GetDBConfig()

	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
	// 	cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	// )

	var err error
	Db, err := sql.Open("mysql", app.Config.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	// Test the connection
	if err := Db.Ping(); err != nil {
		log.Fatalf("Database unreachable: %v", err)
	}
	app.Db = Db

	log.Println("Connected to MySQL!")
}
