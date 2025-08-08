package handlers

import (
	"encoding/json"
	"net/http"
	// "path/filepath"
	// "github.com/go-squad-5/quiz-master/internal/app"
)

type message struct {
	Message string `json:"msg"`
}

func GetQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := message{
		"Hello world!",
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func StoreQuiz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// func (h *Handler) ServeOpenAPISpec(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, filepath.Join("docs", "openapi.yaml"))
// }
//
// func (h *Handler) SwaggerUIHandler() http.Handler {
// 	return http.StripPrefix("/docs/", http.FileServer(http.Dir("docs/swagger")))
// }
//
// func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("OK"))
// }
