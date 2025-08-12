package models

type Question struct {
	Id       string   `json:"ques_id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	Answer   string   `json:"answer"`
}

type Quizzes struct {
	Id         int    `json:"id"`
	SessionID  string `json:"session_id"`
	QuestionID string `json:"question_id"`
	Answer     string `json:"answer"`
}

type CreateQuizBody struct {
	Ssid  string `json:"ssid"`
	Topic string `json:"topic"`
}

type ScoreQuizBody struct {
	Ssid    string `json:"ssid"`
	Answers []struct {
		Id     string `json:"ques_id"`
		Answer string `json:"answer"`
	} `json:"answers"`
}

type ScoreResponse struct {
	Ssid  string `json:"ssid"`
	Score int    `json:"score"`
}
