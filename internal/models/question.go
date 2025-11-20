package models

import "time"

type Question struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Text      string    `gorm:"type:text;not null" json:"text"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	Answers   []Answer  `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"answers,omitempty"`
}
