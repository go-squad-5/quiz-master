package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/go-squad-5/quiz-master/internal/config"
	"github.com/go-squad-5/quiz-master/internal/repositories"
)

type App struct {
	Config      *config.Config
	Repository  *repositories.Repository // and more fields... db, logger, etc.
	ConnChannel chan net.Conn
}

func (app *App) serve() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Failed to accept:", err)
			continue
		}

		app.ConnChannel <- conn
	}
}

func handleConn(connChannel <-chan net.Conn) {
	for conn := range connChannel {

		reader := bufio.NewReader(conn)

		req, err := http.ReadRequest(reader)
		if err != nil {
			fmt.Println("Failed to read request:", err)
			conn.Close()
			return
		}
		rw := newRW(conn)
		// TODO: need to serve request using router
		rw.Flush()
		conn.Close()
	}
}

type rw struct {
	resp        *http.Response
	conn        net.Conn
	body        *bytes.Buffer
	wroteHeader bool
}

func newRW(conn net.Conn) *rw {
	return &rw{
		conn: conn,
		body: new(bytes.Buffer),
		resp: &http.Response{
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Header:     make(http.Header),
			StatusCode: http.StatusOK,
		},
	}
}

func (rw *rw) Header() http.Header {
	return rw.resp.Header
}

func (rw *rw) WriteHeader(statusCode int) {
	if rw.wroteHeader {
		return
	}
	rw.resp.StatusCode = statusCode
	rw.wroteHeader = true
}

func (rw *rw) Write(data []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.body.Write(data)
}

func (rw *rw) Flush() error {
	rw.resp.Body = io.NopCloser(rw.body)
	rw.resp.ContentLength = int64(rw.body.Len())
	return rw.resp.Write(rw.conn)
}

func NewApp() *App {
	config := config.Load()
	repository := InitDB(config.DSN) // InitDB(app)
	connChannel := intializeWorkers(config.WorkerCount)
	log.Println("Connected to MySQL!")
	return &App{
		Config:      config,
		Repository:  repository,
		connChannel: connChannel,
		// Initialize other fields like db, logger, etc.
	}
}

func InitDB(dsn string) *repositories.Repository {
	respository, err := repositories.NewMySQLRepository(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	return &respository
}

func intializeWorkers(workerCount int) chan net.Conn {
	connChannel := make(chan net.Conn)
	for range workerCount {
		go handleConn(connChannel)
	}
	return connChannel
}
