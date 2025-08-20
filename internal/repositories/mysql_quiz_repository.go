package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-squad-5/quiz-master/internal/models"
)

type quizRepository struct {
	db *sql.DB
}

func NewMySQLRepository(dsn string) (Repository, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}

	fmt.Println("Connected to Db!")

	return &quizRepository{
		db: db,
	}, nil
}

func (r *quizRepository) GetAllQuestionByTopic(topic string) ([]models.Question, error) {
	rows, err := r.db.Query("SELECT id, question, options FROM questions WHERE topic = ?", topic)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var quiz []models.Question

	for rows.Next() {
		var optionsJSON []byte
		var q models.Question
		err := rows.Scan(&q.Id, &q.Question, &optionsJSON)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(optionsJSON, &q.Options); err != nil {
			return nil, fmt.Errorf("failed to unmarshal options: %w", err)
		}

		quiz = append(quiz, q)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %s", err)
	}

	return quiz, nil
}

func (r *quizRepository) GetQuestionsByIds(ids []string) ([]models.Question, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("no question IDs provided")
	}

	placeholders := strings.Repeat("?,", len(ids))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf("SELECT id, answer FROM questions WHERE id IN (%s)", placeholders)

	args := make([]any, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var q models.Question
		if err := rows.Scan(&q.Id, &q.Answer); err != nil {
			return nil, fmt.Errorf("failed to scan question: %w", err)
		}
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row error: %w", err)
	}

	return questions, nil
}

func (r *quizRepository) CreateQuiz(ssid string, questions []models.Question) error {
	for _, question := range questions {
		_, err := r.db.Exec(
			"INSERT INTO quizzes (session_id, question_id) VALUES (?, ?)",
			ssid,
			question.Id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *quizRepository) StoreAnswers(ssid, questionId, answer string, is_correct bool) error {
	query := `UPDATE quizzes SET answer = ?, is_correct = ? WHERE session_id = ? and question_id = ?`

	_, err := r.db.Exec(query, answer, is_correct, ssid, questionId)
	if err != nil {
		return fmt.Errorf("failed to store answer: %w", err)
	}
	return nil
}
