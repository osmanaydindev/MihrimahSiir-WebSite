package models

type Author struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Bio         string `json:"bio" gorm:"type:text"`
	BirthYear   *int   `json:"birth_year"`
	DeathYear   *int   `json:"death_year"`
	Nationality string `json:"nationality"`
	Image       string `json:"image"`
	Slug        string `json:"slug" gorm:"unique;not null"`
	IsDeleted   bool   `json:"is_deleted" gorm:"default:false"`
	CreatedAt   string `json:"created_at"`

	// Relationships
	Poems []Poem `json:"poems,omitempty" gorm:"foreignKey:AuthorID"`
	Books []Book `json:"books,omitempty" gorm:"foreignKey:AuthorID"`
}
