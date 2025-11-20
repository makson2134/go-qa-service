package repository

import "github.com/makson2134/go-qa-service/internal/models"

type QuestionRepository interface {
	Create(text string) (*models.Question, error)
	GetByID(id int) (*models.Question, error)
	List() ([]models.Question, error)
	Delete(id int) error
}

type AnswerRepository interface {
	CreateAnswer(questionID int, userID, text string) (*models.Answer, error)
	GetAnswerByID(id int) (*models.Answer, error)
	DeleteAnswer(id int) error
}
