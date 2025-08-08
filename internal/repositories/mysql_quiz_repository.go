package repositories

import "database/sql"

type quizRespository struct {
	db *sql.DB
}

func NewQuizRepository(db *sql.DB) *quizRespository {
	return &quizRespository{
		db: db,
	}
}

// GetAllQuestionByTopic(topic string)
// GetQuestionsByIds(id []string)

func (r *quizRespository) GetAllQuestionByTopic(topic string) {
	err := r.db.QueryRow("SELECT id, question, options from questions WHERE topic = ?", topic).
		Scan(&q.id, &q.question, q.options)
	if err != nil {
		return nil, err
	}
	return &q, nil
}
