package handlers

import (
	"log/slog"

	"github.com/makson2134/go-qa-service/internal/repository"
)

type Handlers struct {
	questions repository.QuestionRepository
	answers   repository.AnswerRepository
	log       *slog.Logger
}

func New(questions repository.QuestionRepository, answers repository.AnswerRepository, log *slog.Logger) *Handlers {
	return &Handlers{
		questions: questions,
		answers:   answers,
		log:       log,
	}
}
