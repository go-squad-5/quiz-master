package repositories

type QuizRepository interface {
	GetAllQuestionByTopic(topic string)
	GetQuestionsByIds(id []string)
}
