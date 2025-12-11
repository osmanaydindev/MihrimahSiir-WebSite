package models

import (
	"time"
)

type Reminder struct {
	ID         uint      `json:"id" autoIncrement:"true"`
	Text       string    `json:"text" bson:"text" binding:"required"`
	Permission int       `json:"permission" bson:"permission"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}
