package postgres

import "github.com/makson2134/go-qa-service/internal/models"

func (db *DB) CreateAnswer(questionID int, userID, text string) (*models.Answer, error) {
	answer := &models.Answer{
		QuestionID: questionID,
		UserID:     userID,
		Text:       text,
	}

	if err := db.conn.Create(answer).Error; err != nil {
		return nil, err
	}

	return answer, nil
}

func (db *DB) GetAnswerByID(id int) (*models.Answer, error) {
	var answer models.Answer

	if err := db.conn.First(&answer, id).Error; err != nil {
		return nil, err
	}

	return &answer, nil
}

func (db *DB) DeleteAnswer(id int) error {
	return db.conn.Delete(&models.Answer{}, id).Error
}
