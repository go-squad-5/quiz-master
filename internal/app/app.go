package app

import (
	"log"
	"net"
	"net/http"

	"github.com/go-squad-5/quiz-master/internal/config"
	"github.com/go-squad-5/quiz-master/internal/repositories"
)

type App struct {
	Config      *config.Config
	Repository  repositories.Repository // and more fields... db, logger, etc.
	ConnChannel chan net.Conn
	Router      http.Handler
}

func (app *App) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept:", err)
			continue
		}

		app.ConnChannel <- conn
	}
}

func getListener(port string) net.Listener {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("can't connect to port", port, ", err: %s", err.Error())
	}
	log.Println("Listening on port: ", port)
	return ln
}

func IntializeWorkers(
	workerCount int,
	worker func(<-chan net.Conn, int),
) chan net.Conn {
	connChannel := make(chan net.Conn)
	for i := range workerCount {
		go worker(connChannel, i)
		// time.Sleep(500 * time.Millisecond)
	}
	return connChannel
}
