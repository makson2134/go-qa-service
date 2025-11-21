package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/makson2134/go-qa-service/internal/api/dto"
	"gorm.io/gorm"
)

func (h *Handlers) ListQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.questions.List()
	if err != nil {
		h.log.Error("failed to list questions", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	response := make([]dto.QuestionResponse, len(questions))
	for i, q := range questions {
		response[i] = dto.QuestionResponse{
			ID:        q.ID,
			Text:      q.Text,
			CreatedAt: q.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

func (h *Handlers) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "Text cannot be empty", http.StatusBadRequest)
		return
	}

	question, err := h.questions.Create(req.Text)
	if err != nil {
		h.log.Error("failed to create question", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	response := dto.QuestionResponse{
		ID:        question.ID,
		Text:      question.Text,
		CreatedAt: question.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

func (h *Handlers) GetQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/questions/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := h.questions.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Question not found", http.StatusNotFound)

			return
		}

		h.log.Error("failed to get question", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	answers := make([]dto.AnswerResponse, len(question.Answers))
	for i, a := range question.Answers {
		answers[i] = dto.AnswerResponse{
			ID:         a.ID,
			QuestionID: a.QuestionID,
			UserID:     a.UserID,
			Text:       a.Text,
			CreatedAt:  a.CreatedAt,
		}
	}

	response := dto.QuestionWithAnswersResponse{
		ID:        question.ID,
		Text:      question.Text,
		CreatedAt: question.CreatedAt,
		Answers:   answers,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

func (h *Handlers) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/questions/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := h.questions.Delete(id); err != nil {
		h.log.Error("failed to delete question", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
