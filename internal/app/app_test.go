package app

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/go-squad-5/quiz-master/internal/config"
	// "github.com/go-squad-5/quiz-master/internal/repositories"
)

func Test_Serve(t *testing.T) {
	app := App{
		Config: &config.Config{
			Port:        ":8090",
			DSN:         "jjj",
			WorkerCount: 9,
		},
		// Repository:  repositories.Repository,
		ConnChannel: make(chan net.Conn),
	}

	go func() {
		err := app.Serve()
		if err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()
	// waiting for server to start
	time.Sleep(1 * time.Second)

	conn, err := net.Dial("tcp", fmt.Sprintf("localhost%s", app.Config.Port))
	if err != nil {
		t.Errorf("Sending connection to server and failed.")
		return
	}
	defer conn.Close()

	select {
	case recieveSentRequest := <-app.ConnChannel:
		if recieveSentRequest == nil {
			t.Errorf("when recieving connection from server, recieved a nil connection")
		}
	case <-time.After(3 * time.Second):
		t.Error("when recieving connection from server, timeout no connection recieved")
	}
}

func Test_intializeWorkers(t *testing.T) {
	log.SetOutput(io.Discard)
	// req := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	// router := http.NewServeMux()
	//
	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("request received"))
	// })

	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "test sending 1 request",
			count: 1,
		},
		{
			name:  "test sending 2 request",
			count: 2,
		},
		{
			name:  "test sending 3 request",
			count: 3,
		},
	}

	ConnChannel := IntializeWorkers(2, worker)
	time.Sleep(1 * time.Second)

	for _, test := range tests {

		clientConn := make([]net.Conn, test.count)
		serverConn := make([]net.Conn, test.count)

		for i := range test.count {
			clientConn[i], serverConn[i] = net.Pipe()
		}

		// var wg sync.WaitGroup
		for i := range test.count {
			go func(i int) {
				ConnChannel <- serverConn[i]
				time.Sleep(time.Millisecond)

				clientConn[i].Write([]byte("send request"))
				// _, err := clientConn[i].Write([]byte("send request"))
				// if err != nil {
				// 	t.Errorf("Failed to write to reqeust to connection for %s", test.name)
				// }
			}(i)
		}

		for i := range test.count {

			reader := bufio.NewReader(clientConn[i])
			resp := make([]byte, 16)
			_, err := reader.Read(resp)
			if err != nil {
				t.Errorf("failed to read response for %s", test.name)
			}
			if string(resp) != "request received" {
				t.Errorf(
					"for %s Expected response to be %s got %s",
					test.name,
					"request received",
					string(resp),
				)
			}

			// if err != nil {
			// 	t.Fatalf("Failed to parse response: %v", err)
			// }
			// defer resp.Body.Close()

			// body, err := io.ReadAll(resp.Body)
			// if err != nil {
			// 	t.Fatalf("Failed to read response body: %v", err)
			// }
			//
			// if string(body) != "request received" {
			// 	t.Errorf("Expected body to be %q, got %q", "request received", string(body))
			// }

		}

		for i := range test.count {
			serverConn[i].Close()
			clientConn[i].Close()
		}
	}
}

func worker(channel <-chan net.Conn, id int) {
	for conn := range channel {

		conn.Write([]byte("request received"))
		conn.Close()

	}
}
