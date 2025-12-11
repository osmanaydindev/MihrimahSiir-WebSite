package models

type Comment struct {
	ID        uint   `json:"id" autoIncrement:"true"`
	AdminID   uint   `json:"admin_id"`
	Admin     Admin  `json:"admin" gorm:"foreignKey:AdminID"`
	BookID    uint   `json:"book_id"`
	Book      Book   `json:"book" gorm:"foreignKey:BookID"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Page      *int   `json:"page"` // Optional page number
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`
}
