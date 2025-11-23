package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/makson2134/go-qa-service/internal/models"
	"github.com/makson2134/go-qa-service/pkg"
	"gorm.io/gorm"
)

func TestCreateAnswer_QuestionNotFound(t *testing.T) {
	mockQuestions := &mockQuestionRepo{
		getByIDFunc: func(id int) (*models.Question, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}

	logger := pkg.NewLogger("error", "json")
	h := New(mockQuestions, &mockAnswerRepo{}, logger)

	body := map[string]string{
		"user_id": "user-123",
		"text":    "Some answer",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/questions/999/answers/", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	h.CreateAnswer(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestCreateAnswer_EmptyText(t *testing.T) {
	logger := pkg.NewLogger("error", "json")
	h := New(&mockQuestionRepo{}, &mockAnswerRepo{}, logger)

	tests := []struct {
		name string
		body map[string]string
	}{
		{
			name: "empty string",
			body: map[string]string{
				"user_id": "user-123",
				"text":    "",
			},
		},
		{
			name: "whitespace only",
			body: map[string]string{
				"user_id": "user-123",
				"text":    "   ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			h.CreateAnswer(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}

func TestCreateAnswer_EmptyUserID(t *testing.T) {
	logger := pkg.NewLogger("error", "json")
	h := New(&mockQuestionRepo{}, &mockAnswerRepo{}, logger)

	tests := []struct {
		name string
		body map[string]string
	}{
		{
			name: "empty string",
			body: map[string]string{
				"user_id": "",
				"text":    "Some answer",
			},
		},
		{
			name: "whitespace only",
			body: map[string]string{
				"user_id": "   ",
				"text":    "Some answer",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/questions/1/answers/", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			h.CreateAnswer(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}
