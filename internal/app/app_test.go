package app

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-squad-5/quiz-master/internal/config"
)

type mockListener struct {
	acceptSendError bool
	calledAccept    bool
}

func (m *mockListener) Accept() (net.Conn, error) {
	if m.calledAccept {
		// time.Sleep(1 * time.Minute)
		wg := sync.WaitGroup{}
		wg.Add(1)
		wg.Wait()
	}
	m.calledAccept = true
	if m.acceptSendError {
		return nil, errors.New("got error calling accept")
	}

	return &net.TCPConn{}, nil
}

func (m mockListener) Listen() (net.Listener, error) {
	return nil, nil
}

func (m mockListener) Addr() net.Addr {
	return nil
}

func (m mockListener) Close() error {
	return nil
}

func Test_Serve_AcceptWithoutError(t *testing.T) {
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
		err := app.Serve(&mockListener{acceptSendError: false, calledAccept: false})
		if err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()
	// waiting for server to start
	time.Sleep(1 * time.Second)

	// conn, err := net.Dial("tcp", fmt.Sprintf("localhost%s", app.Config.Port))
	// if err != nil {
	// 	t.Errorf("Sending connection to server and failed.")
	// 	return
	// }
	// defer conn.Close()

	select {
	case recieveSentRequest := <-app.ConnChannel:
		if recieveSentRequest == nil {
			t.Errorf("when recieving connection from server, recieved a nil connection")
		}
	case <-time.After(3 * time.Second):
		t.Error("when recieving connection from server, timeout no connection recieved")
	}
}

func Test_Serve_AcceptWithError(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	app := App{
		Config: &config.Config{
			Port:        ":8090",
			DSN:         "jjj",
			WorkerCount: 9,
		},
		ConnChannel: make(chan net.Conn),
	}

	go func() {
		err := app.Serve(&mockListener{acceptSendError: true, calledAccept: false})
		if err != nil {
			t.Errorf("Server error: %v", err)
		}
	}()
	// waiting for server to start
	time.Sleep(1 * time.Second)

	log.SetOutput(os.Stderr)
	output := buf.String()
	expected := "got error calling accept"

	if !strings.Contains(output, expected) {
		t.Errorf("expected error: %s, got error: %s.", expected, output)
	}
}

func Test_intializeWorkers(t *testing.T) {
	log.SetOutput(io.Discard)

	tests := []struct {
		name  string
		count int
	}{
		{
			name:  "test sending 1 request, less than workers count",
			count: 1,
		},
		{
			name:  "test sending 2 request, equal to worker count",
			count: 2,
		},
		{
			name:  "test sending 3 request, more than worker count",
			count: 3,
		},
	}

	ConnChannel := IntializeWorkers(2, worker)
	time.Sleep(1 * time.Second)

	for _, test := range tests {
		func() {
			clientConn := make([]net.Conn, test.count)
			serverConn := make([]net.Conn, test.count)

			// intialize connections
			for i := range test.count {
				clientConn[i], serverConn[i] = net.Pipe()
				defer clientConn[i].Close()
			}

			// send and write request to connection
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
