package models

type MihrimahCard struct {
	ID      uint   `json:"id" autoIncrement:"true"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
