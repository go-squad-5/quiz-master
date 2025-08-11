package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-squad-5/quiz-master/internal/models"
	"github.com/go-squad-5/quiz-master/internal/repositories"
	// "path/filepath"
	// "github.com/go-squad-5/quiz-master/internal/app"
)

type message struct {
	Message string `json:"msg"`
}

type Handler struct {
	repo repositories.Repository
}

func NewHandler(repo repositories.Repository) Handler {
	return Handler{
		repo,
	}
}

func (h *Handler) GetQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	topic := r.URL.Query().Get("topic")

	var req models.CreateQuizBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ssid := req.Ssid

	quiz, err := h.repo.GetAllQuestionByTopic(topic)
	if err != nil {
		fmt.Printf("Error getting questions: %v", err)
		// log.Println("Error getting the questions")
	}

	if len(quiz) == 0 {
		http.Error(w, "No questions found", http.StatusNotFound)
	}

	if err := h.repo.CreateQuiz(ssid, quiz); err != nil {
		http.Error(w, "Unable to create quiz", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(quiz); err != nil {
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ids := make([]string, len(req.Questions))

	for _, question := range req.Questions {
		ids = append(ids, question.Id)
	}

	questions, err := h.repo.GetQuestionsByIds(ids)
	if err != nil {
		log.Println(err)
		return
	}

	for _, question := range questions {
		for _, reqQuestion := range req.Questions {
			if reqQuestion.Id == question.Id {
				err := h.repo.StoreAnswers(req.Ssid, question.Id, reqQuestion.Answer)
				if err != nil {
					log.Println("Failed to save answer: ", err)
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

	// var scoreRes models.ScoreResponse
	//
	//  json.Marshal(&scoreRes)
	scoreRes := fmt.Sprintf(`{"score": %d}`, score)
	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte(scoreRes))
}
