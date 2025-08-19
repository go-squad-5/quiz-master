package app

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"testing"
)

func Test_HandleConn(t *testing.T) {
	log.SetOutput(io.Discard)
	var app App

	req := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	connChannel := make(chan net.Conn)
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("request received"))
	})

	app.Router = router
	clientConn, ServerConn := net.Pipe()

	go app.HandleConn(connChannel, 0)

	connChannel <- ServerConn

	clientConn.Write([]byte(req))

	reader := bufio.NewReader(clientConn)
	resp, err := http.ReadResponse(reader, nil)
	if err != nil {
		t.Errorf("failed to read response")
	}
	defer resp.Body.Close()
	body := make([]byte, 1024)
	n, _ := resp.Body.Read(body)

	if string(body[:n]) != "request received" {
		t.Errorf("invalid response expected: %s got: %s", "request received", string(body[:n]))
	}
}
