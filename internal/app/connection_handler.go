package app

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func (app *App) HandleConn(connChannel <-chan net.Conn, id int) {
	log.Println("handleConn ", id, " started")
	requestNum := 1
	for conn := range connChannel {

		reader := bufio.NewReader(conn)

		rw := newRW(conn)
		req, err := http.ReadRequest(reader)
		if err != nil {
			http.Error(rw, "Invalid request", http.StatusBadRequest)
			rw.Flush()
			conn.Close()
			return
		}
		log.Printf(
			"routine: %d processing request\t\t%s\trequest id: %d_%d",
			id,
			req.URL,
			id,
			requestNum,
		)
		time.Sleep(time.Duration(rand.Intn(6)+1) * time.Second)
		app.Router.ServeHTTP(rw, req)
		log.Printf(
			"routine: %d sending response for request\t%s\trequest id: %d_%d",
			id,
			req.URL,
			id,
			requestNum,
		)
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
