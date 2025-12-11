package models

type Book struct {
	ID        uint   `json:"id" autoIncrement:"true"`
	Name      string `json:"name"`
	Author    string `json:"author"` // Deprecated: kept for backward compatibility
	AuthorID  *uint  `json:"author_id"`
	Slug      string `json:"slug"`
	Image     string `json:"image"`
	Page      int    `json:"page"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false"`
	CreatedAt string `json:"created_at"`
	Community int    `json:"community" gorm:"default:1"` // 1=private (role_id 1,2), 2=public (all)

	// Relationships
	Comments   []Comment `json:"comments" gorm:"foreignKey:BookID"`
	AuthorData *Author   `json:"author_data,omitempty" gorm:"foreignKey:AuthorID"`
}
