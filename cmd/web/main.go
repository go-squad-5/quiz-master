package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-squad-5/quiz-master/internal/app"
	"github.com/go-squad-5/quiz-master/internal/config"
	"github.com/go-squad-5/quiz-master/internal/handlers"
	"github.com/go-squad-5/quiz-master/internal/repositories"
	"github.com/go-squad-5/quiz-master/internal/router"
)

func main() {
	var application app.App
	config := config.LoadConfig()

	respository, err := repositories.NewMySQLRepository(config.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	log.Println("Connected to MySQL!")

	handler := handlers.NewHandler(respository)
	router := router.NewRouter(handler)
	connChannel := app.IntializeWorkers(config.WorkerCount, application.HandleConn)

	application.Config = config
	application.Repository = respository
	application.ConnChannel = connChannel
	application.Router = router

	log.Fatal(application.Serve(app.GetListener(application.Config.Port)))
}
