package models

import "time"

type Answer struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`
	QuestionID int       `gorm:"not null;index" json:"question_id"`
	UserID     string    `gorm:"type:varchar(255);not null" json:"user_id"`
	Text       string    `gorm:"type:text;not null" json:"text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
