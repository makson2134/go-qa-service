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

func (h *Handlers) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 2 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	questionID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req dto.CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.UserID) == "" {
		http.Error(w, "UserID cannot be empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Text) == "" {
		http.Error(w, "Text cannot be empty", http.StatusBadRequest)
		return
	}

	_, err = h.questions.GetByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Question not found", http.StatusNotFound)
			return
		}

		h.log.Error("failed to check question existence", "error", err, "question_id", questionID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	answer, err := h.answers.CreateAnswer(questionID, req.UserID, req.Text)
	if err != nil {
		h.log.Error("failed to create answer", "error", err, "question_id", questionID)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	response := dto.AnswerResponse{
		ID:         answer.ID,
		QuestionID: answer.QuestionID,
		UserID:     answer.UserID,
		Text:       answer.Text,
		CreatedAt:  answer.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

func (h *Handlers) GetAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/answers/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	answer, err := h.answers.GetAnswerByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Answer not found", http.StatusNotFound)
			return
		}

		h.log.Error("failed to get answer", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	response := dto.AnswerResponse{
		ID:         answer.ID,
		QuestionID: answer.QuestionID,
		UserID:     answer.UserID,
		Text:       answer.Text,
		CreatedAt:  answer.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.log.Error("failed to encode response", "error", err)
	}
}

func (h *Handlers) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/answers/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid answer ID", http.StatusBadRequest)
		return
	}

	if err := h.answers.DeleteAnswer(id); err != nil {
		h.log.Error("failed to delete answer", "error", err, "id", id)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
