package app

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"testing"
)

func Test_HandleConn(t *testing.T) {
	// log.SetOutput(io.Discard)
	var app App

	tests := []struct {
		name string
		req  string
		res  string
	}{
		{
			name: "valid request",
			req:  "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n",
			res:  "response send",
		},
		{
			name: "invalid request",
			req:  "invalid request\r\n\r\n\r\n",
			res:  "Invalid request\n",
		},
	}
	connChannel := make(chan net.Conn)
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("response send"))
	})

	app.Router = router

	go app.HandleConn(connChannel, 0)

	for _, test := range tests {
		clientConn, ServerConn := net.Pipe()
		connChannel <- ServerConn

		_, err := clientConn.Write([]byte(test.req))
		if err != nil {
			log.Println("failed to write request to clientConn: ", err)
		}

		reader := bufio.NewReader(clientConn)
		resp, err := http.ReadResponse(reader, nil)
		if err != nil {
			t.Errorf("failed to read response")
		}
		defer resp.Body.Close()
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)

		if string(body[:n]) != test.res {
			t.Errorf("for %s expected: %s got: %s", test.name, "request received", string(body[:n]))
		}
		clientConn.Close()
	}
}

func Test_WrittenHeaderTwice(t *testing.T) {
	clientConn, ServerConn := net.Pipe()
	defer clientConn.Close()
	defer ServerConn.Close()
	rw := newRW(ServerConn)

	rw.WriteHeader(http.StatusBadRequest)
	rw.WriteHeader(http.StatusOK)

	if rw.resp.StatusCode == http.StatusOK {
		t.Errorf("response writter wrote headers twice")
	}
}
