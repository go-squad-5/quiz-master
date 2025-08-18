package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"

	"github.com/go-squad-5/quiz-master/internal/router"
)

func handleConn(connChannel <-chan net.Conn, id int, router *router.Router) {
	log.Println("handleConn ", id, " started")
	for conn := range connChannel {

		reader := bufio.NewReader(conn)

		req, err := http.ReadRequest(reader)
		if err != nil {
			fmt.Println("Failed to read request:", err)
			conn.Close()
			return
		}
		rw := newRW(conn)
		log.Println("routine: ", id, "processing request")
		time.Sleep(time.Duration(rand.Intn(6)+1) * time.Second)
		router.ServeHTTP(rw, req)
		log.Println("routine: ", id, "sending response")
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
