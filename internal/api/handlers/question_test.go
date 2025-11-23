package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/makson2134/go-qa-service/internal/models"
	"github.com/makson2134/go-qa-service/pkg"
)

type mockQuestionRepo struct{}

func (m *mockQuestionRepo) Create(text string) (*models.Question, error) {
	return nil, nil
}

func (m *mockQuestionRepo) GetByID(id int) (*models.Question, error) {
	return nil, nil
}

func (m *mockQuestionRepo) List() ([]models.Question, error) {
	return nil, nil
}

func (m *mockQuestionRepo) Delete(id int) error {
	return nil
}

type mockAnswerRepo struct{}

func (m *mockAnswerRepo) CreateAnswer(questionID int, userID, text string) (*models.Answer, error) {
	return nil, nil
}

func (m *mockAnswerRepo) GetAnswerByID(id int) (*models.Answer, error) {
	return nil, nil
}

func (m *mockAnswerRepo) DeleteAnswer(id int) error {
	return nil
}

func TestCreateQuestion_EmptyText(t *testing.T) {
	logger := pkg.NewLogger("error", "json")
	h := New(&mockQuestionRepo{}, &mockAnswerRepo{}, logger)

	tests := []struct {
		name string
		body map[string]string
	}{
		{
			name: "empty string",
			body: map[string]string{"text": ""},
		},
		{
			name: "whitespace only",
			body: map[string]string{"text": "   "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			h.CreateQuestion(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
			}
		})
	}
}
