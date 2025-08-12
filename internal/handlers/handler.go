package handlers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/go-squad-5/quiz-master/internal/models"
	"github.com/go-squad-5/quiz-master/internal/repositories"
)

type Handler struct {
	repo repositories.Repository
}

func NewHandler(repo repositories.Repository) Handler {
	return Handler{
		repo,
	}
}

type ResponseGetQuiz struct {
	Questions []models.Question `json:"questions"`
}

func (h *Handler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		// log.Println("Request recieved at invalid method want: ", http.MethodPost, "got: ", r.Method)
		return
	}

	var req models.CreateQuizBody

	noOfQuestions, err := strconv.Atoi(r.URL.Query().Get("noOfQuestions"))
	if err != nil {
		noOfQuestions = 20
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// log.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ssid := req.Ssid

	quiz, err := h.repo.GetAllQuestionByTopic(req.Topic)
	if err != nil {
		http.Error(w, "Error getting topic", http.StatusInternalServerError)
		// log.Printf("Error getting questions: %v", err)
	}

	if len(quiz) < noOfQuestions {
		// log.Println("No questions found for the given topic")
		http.Error(w, "Not enough questions found", http.StatusNotFound)
	}
	randomNumber := rand.Intn(len(quiz) - noOfQuestions)
	quiz = quiz[randomNumber : noOfQuestions+randomNumber]

	if err := h.repo.CreateQuiz(ssid, quiz); err != nil {
		http.Error(w, "Unable to create quiz", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	quizResponse := ResponseGetQuiz{
		Questions: quiz,
	}

	if err := json.NewEncoder(w).Encode(quizResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *Handler) ScoreQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ScoreQuizBody
	var score int

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// log.Printf("error: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ids := make([]string, 0, len(req.Answers))

	for _, answer := range req.Answers {
		ids = append(ids, answer.Id)
	}

	questions, err := h.repo.GetQuestionsByIds(ids)
	if err != nil {
		// log.Println(err)
		http.Error(w, "Error verifying questions", http.StatusInternalServerError)
		return
	}

	for _, question := range questions {
		for _, reqQuestion := range req.Answers {
			if reqQuestion.Id == question.Id {
				channel := make(chan error)
				go func(
					ssid,
					quesionId,
					answer string,
					isCorrect bool,
					channel chan<- error,
				) {
					channel <- h.repo.StoreAnswers(ssid, quesionId, answer, isCorrect)
				}(req.Ssid, question.Id, reqQuestion.Answer, reqQuestion.Answer == question.Answer, channel)
				err := <-channel
				if err != nil {
					// log.Println("Failed to save answer: ", err)
					http.Error(w, "Failed to save answer", http.StatusInternalServerError)
					return
				}
				if reqQuestion.Answer == question.Answer {
					score++
					break
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.ScoreResponse{
		Ssid:  req.Ssid,
		Score: score,
	})
}
