package dto

import "time"

type CreateQuestionRequest struct {
	Text string `json:"text"`
}

type QuestionResponse struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type QuestionWithAnswersResponse struct {
	ID        int              `json:"id"`
	Text      string           `json:"text"`
	CreatedAt time.Time        `json:"created_at"`
	Answers   []AnswerResponse `json:"answers"`
}
