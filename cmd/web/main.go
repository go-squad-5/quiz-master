package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-squad-5/quiz-master/internal/app"
	"github.com/go-squad-5/quiz-master/internal/router"
)

func main() {
	app := app.NewApp()

	router := router.NewRouter(app)
	// route := router.Route()

	log.Println("Starting server on", app.Config.Port)
	log.Fatal(http.ListenAndServe(app.Config.Port, router))
}
