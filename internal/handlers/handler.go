package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		return
	}
	// topic := r.URL.Query().Get("topic")

	var req models.CreateQuizBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ssid := req.Ssid
	log.Println("Session ID:", ssid)
	log.Println("Topic:", req.Topic)

	quiz, err := h.repo.GetAllQuestionByTopic(req.Topic)
	if err != nil {
		http.Error(w, "Error getting topic", http.StatusInternalServerError)
		fmt.Printf("Error getting questions: %v", err)
		// log.Println("Error getting the questions")
	}

	if len(quiz) == 0 {
		log.Println("No questions found for the given topic")
		http.Error(w, "No questions found", http.StatusNotFound)
	}

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
		log.Printf("error: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ids := make([]string, 0, len(req.Answers))

	for _, answer := range req.Answers {
		ids = append(ids, answer.Id)
	}

	fmt.Println(ids)

	questions, err := h.repo.GetQuestionsByIds(ids)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error verifying questions", http.StatusInternalServerError)
		return
	}

	for _, question := range questions {
		for _, reqQuestion := range req.Answers {
			if reqQuestion.Id == question.Id {
				err := h.repo.StoreAnswers(
					req.Ssid,
					question.Id,
					reqQuestion.Answer,
					reqQuestion.Answer == question.Answer,
				)
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.ScoreResponse{
		Ssid:  req.Ssid,
		Score: score,
	})
}
