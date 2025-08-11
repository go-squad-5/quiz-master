package models

// package config

type Question struct {
	Id       string   `json:"id"`
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
	Ssid string `json:"ssid"`
}

type ScoreQuizBody struct {
	Ssid      string `json:"ssid"`
	Questions []struct {
		Id     string `json:"id"`
		Answer string `json:"answer"`
	}
}

type ScoreResponse struct {
	Score int `json:"score"`
}
