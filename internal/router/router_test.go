package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Mockhandler struct{}

func (h *Mockhandler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("send response by GetQuiz"))
}

func (h *Mockhandler) ScoreQuiz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("send response by ScoreQuiz"))
}

func Test_serveHTTP(t *testing.T) {
	mockhandler := &Mockhandler{}
	router := NewRouter(mockhandler)

	tests := []struct {
		name             string
		url              string
		exceptedResponse string
	}{
		{
			name:             "testing router.serveHTTP sending url /quiz/fetch",
			url:              "/quiz/fetch",
			exceptedResponse: "send response by GetQuiz",
		},
		{
			name:             "testing router.serveHTTP sending url /quiz/score",
			url:              "/quiz/score",
			exceptedResponse: "send response by ScoreQuiz",
		},
		{
			name:             "testing router.serveHTTP sending url /fadls;hf",
			url:              "/fadls;hf",
			exceptedResponse: "404 page not found\n",
		},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", test.url, nil)
		router.ServeHTTP(rr, req)
		body := rr.Body.String()
		if body != test.exceptedResponse {
			t.Errorf("%s excepted: %s, got: %s.", test.name, test.exceptedResponse, body)
		}
	}
}
