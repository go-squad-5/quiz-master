package app

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
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
	// log.SetOutput(io.Discard)

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
		func() {
			clientConn := make([]net.Conn, test.count)
			serverConn := make([]net.Conn, test.count)

			for i := range test.count {
				clientConn[i], serverConn[i] = net.Pipe()
				defer clientConn[i].Close()
			}

			wg := sync.WaitGroup{}
			for i := range test.count {
				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					ConnChannel <- serverConn[i]

					_, err := clientConn[i].Write([]byte("send request\r\n"))
					if err != nil {
						t.Errorf("Failed to write to reqeust to connection: %s", err.Error())
					}
				}(i)
			}

			for i := range test.count {

				reader := bufio.NewReader(clientConn[i])
				resp := make([]byte, 50)
				n, err := reader.Read(resp)
				resp = resp[:n]
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
			}
			wg.Wait()
		}()
	}
}

func worker(channel <-chan net.Conn, id int) {
	for conn := range channel {
		resp := make([]byte, 50)
		_, err := conn.Read(resp)
		if err != nil {
			log.Println("%w", err)
		}

		conn.Write([]byte("request received"))
		time.Sleep(1 * time.Second)
		conn.Close()

	}
}
