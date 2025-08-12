package main

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-squad-5/quiz-master/internal/app"
)

func main() {
	app := app.NewApp()

	log.Fatal(app.Serve())
}
