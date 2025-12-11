package models

type Poem struct {
	ID             uint   `json:"id" autoIncrement:"true"`
	Title          string `json:"title"`
	Author         string `json:"author"` // Deprecated: kept for backward compatibility
	AuthorID       *uint  `json:"author_id"`
	Content        string `json:"content"`
	IsDeleted      bool   `json:"is_deleted" gorm:"default:false"`
	Slug           string `json:"slug"`
	CreatedAt      string `json:"created_at"`
	CreatedAtParse string `json:"created_at_parse"`
	Community      int    `json:"community" gorm:"default:1"` // 1=private (role_id 1,2), 2=public (role_id 3)
	LikeCount      int    `json:"like_count" gorm:"-"`        // Computed field, not stored in DB

	// Relationship
	AuthorData *Author `json:"author_data,omitempty" gorm:"foreignKey:AuthorID"`
}
