package models

type Homepage struct {
	ID         uint   `json:"id" autoIncrement:"true"`
	Title      string `json:"title"`
	Subtitle   string `json:"subtitle"`
	Content    string `json:"content"`
	Permission int    `json:"permission"` // 1=Admin, 2=Member, 3=Guest
}
