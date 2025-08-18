package app

import (
	"log"
	"net"
	"time"

	"github.com/go-squad-5/quiz-master/internal/config"
	"github.com/go-squad-5/quiz-master/internal/repositories"
	"github.com/go-squad-5/quiz-master/internal/router"
)

type App struct {
	Config      *config.Config
	Repository  repositories.Repository // and more fields... db, logger, etc.
	ConnChannel chan net.Conn
	router      *router.Router
}

func (app *App) Serve() error {
	ln, err := net.Listen("tcp", app.Config.Port)
	if err != nil {
		return err
	}
	log.Println("Listening on port: ", app.Config.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept:", err)
			continue
		}

		app.ConnChannel <- conn
	}
}

func IntializeWorkers(workerCount int, router *router.Router) chan net.Conn {
	connChannel := make(chan net.Conn)
	for i := range workerCount {
		go handleConn(connChannel, i, router)
		time.Sleep(500 * time.Millisecond)
	}
	return connChannel
}
