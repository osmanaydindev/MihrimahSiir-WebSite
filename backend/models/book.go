package models

type Book struct {
	ID          uint   `json:"id" autoIncrement:"true"`
	Name        string `json:"name"`
	Author      string `json:"author"` // Deprecated: kept for backward compatibility
	AuthorID    *uint  `json:"author_id"`
	Slug        string `json:"slug"`
	Image       string `json:"image"`
	Page        int    `json:"page"`
	ISBN        string `json:"isbn" gorm:"type:varchar(13);index"` // normalize edilmiş ISBN-13
	Description string `json:"description" gorm:"type:text"`
	IsDeleted   bool   `json:"is_deleted" gorm:"default:false"`
	CreatedAt   string `json:"created_at"`
	// 1=özel (role_id 1,2), 2=herkese açık, 3=sadece book_visibilities'teki kullanıcılar
	Community int `json:"community" gorm:"default:1"`

	// Relationships
	Comments   []Comment `json:"comments" gorm:"foreignKey:BookID"`
	AuthorData *Author   `json:"author_data,omitempty" gorm:"foreignKey:AuthorID"`
}
