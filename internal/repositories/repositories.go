package repositories

import "github.com/go-squad-5/quiz-master/internal/models"

type Repository interface {
	GetAllQuestionByTopic(topic string) ([]models.Question, error)
	GetQuestionsByIds(id []string) ([]models.Question, error)
	CreateQuiz(string, []models.Question) error
	StoreAnswers(string, string, string, bool) error
}
