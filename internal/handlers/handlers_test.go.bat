package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-squad-5/quiz-master/internal/models"
)

// type Repository interface {
// 	GetAllQuestionByTopic(topic string) ([]models.Question, error)
// 	GetQuestionsByIds(id []string) ([]models.Question, error)
// 	CreateQuiz(string, []models.Question) error
// 	StoreAnswers(string, string, string, bool) error
// }

type mockrepo struct{}

func (*mockrepo) GetAllQuestionByTopic(topic string) ([]models.Question, error) {
	if topic == "invalid topic" {
		return nil, fmt.Errorf("")
	}
	return []models.Question{
		{
			Id:       "1",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},

		{
			Id:       "2",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "3",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "4",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "5",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "6",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},

		{
			Id:       "7",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "8",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "9",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "10",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "11",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "12",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "13",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "14",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "15",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "16",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "17",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "18",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "19",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
		{
			Id:       "20",
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		},
	}, nil
}

func (*mockrepo) GetQuestionsByIds(id []string) ([]models.Question, error) {
	result := make([]models.Question, 0)
	for _, v := range id {
		result = append(result, models.Question{
			Id:       v,
			Question: "ques1",
			Options:  []string{"a", "b", "c", "c"},
			Answer:   "",
		})
	}
	return result, nil
}

func (*mockrepo) CreateQuiz(ssid string, questions []models.Question) error {
	if ssid == "invalid" {
		return fmt.Errorf("")
	}
	return nil
}

func (*mockrepo) StoreAnswers(string, string, string, bool) error {
	return nil
}

func Test_GetQuiz(t *testing.T) {
	repo := &mockrepo{}
	h := NewHandler(repo)

	tests := []struct {
		name             string
		method           string
		url              string
		body             io.Reader
		status           int
		execptedResponse string
	}{
		{
			name:             "invalid method",
			method:           "PATCH",
			url:              "/quiz/fetch",
			status:           http.StatusMethodNotAllowed,
			execptedResponse: "Method Not Allowed\n",
		},
		{
			name:             "invalid noOfQuestions",
			method:           http.MethodPost,
			url:              "/quiz/fetch?noOfQuestions=10",
			body:             strings.NewReader(`{"topic": "valid", "ssid": "abc123"}`),
			status:           http.StatusInternalServerError,
			execptedResponse: "Not enough questions found\n",
		},
		{
			name:             "invalid request body",
			method:           http.MethodPost,
			url:              "/quiz/fetch",
			body:             strings.NewReader("invalid-json"),
			status:           http.StatusBadRequest,
			execptedResponse: "Invalid request body\n",
		},
		{
			name:             "invalid topic",
			method:           http.MethodPost,
			body:             strings.NewReader(`{"topic": "invalid topic", "ssid": "abc123"}`),
			url:              "/quiz/fetch",
			status:           http.StatusInternalServerError,
			execptedResponse: "Error getting topic",
		},
		{
			name:             "invalid ssid",
			method:           http.MethodPost,
			url:              "/quiz/fetch?noOfQuestions=2",
			body:             strings.NewReader(`{"topic": "valid", "ssid": "invalid"}`),
			status:           http.StatusInternalServerError,
			execptedResponse: "Unable to create quiz\n",
		},
		// {
		// 	name:             "valid request",
		// 	method:           http.MethodPost,
		// 	url:              "/quiz/fetch?noOfQuestions=2",
		// 	body:             strings.NewReader(`{"topic": "valid", "ssid": "abc123"}`),
		// 	status:           http.StatusOK,
		// 	execptedResponse: json.Unmarshal(data []byte, v any),
		// },
		{},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(test.method, test.url, test.body)

		h.GetQuiz(rr, req)

		// test status code
		if rr.Code != test.status {
			t.Errorf("for %s expected status %d, got %d", test.name, test.status, rr.Code)
		}
		// test body
		body := rr.Body.String()

		if body != test.execptedResponse {
			t.Errorf("for %s, exepected body: %s, got: %s.", test.name, test.body, body)
		}
	}
}

func Test_ScoreQuiz(t *testing.T) {
	repo := &mockrepo{}
	h := NewHandler(repo)

	tests := []struct {
		name             string
		method           string
		url              string
		body             io.Reader
		status           int
		execptedResponse string
	}{}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		req := httptest.NewRequest(test.method, test.url, test.body)

		h.GetQuiz(rr, req)

		// test status code
		if rr.Code != test.status {
			t.Errorf("for %s expected status %d, got %d", test.name, test.status, rr.Code)
		}
		// test body
		body := rr.Body.String()

		if body != test.execptedResponse {
			t.Errorf("for %s, exepected body: %s, got: %s.", test.name, test.body, body)
		}
	}
}
