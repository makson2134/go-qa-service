package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

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

func (h *Handlers) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "ok"}); err != nil {
		h.log.Error("Failed to encode health check response", "error", err)
	}
}
