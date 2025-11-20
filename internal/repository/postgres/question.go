package postgres

import "github.com/makson2134/go-qa-service/internal/models"

func (db *DB) Create(text string) (*models.Question, error) {
	question := &models.Question{Text: text}

	if err := db.conn.Create(question).Error; err != nil {
		return nil, err
	}

	return question, nil
}

func (db *DB) GetByID(id int) (*models.Question, error) {
	var question models.Question

	if err := db.conn.Preload("Answers").First(&question, id).Error; err != nil {
		return nil, err
	}

	return &question, nil
}

func (db *DB) List() ([]models.Question, error) {
	var questions []models.Question

	if err := db.conn.Find(&questions).Error; err != nil {
		return nil, err
	}
	
	return questions, nil
}

func (db *DB) Delete(id int) error {
	return db.conn.Delete(&models.Question{}, id).Error
}
