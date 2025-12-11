package models

type Admin struct {
	ID                 uint      `json:"id" autoIncrement:"true"`
	Username           string    `json:"username" gorm:"unique"`
	Email              string    `json:"email" gorm:"unique"`
	Password           []byte    `json:"password" readOnly:"false"`
	RoleID             uint      `json:"role_id"`
	ProfileImage       string    `json:"profile_image"`
	IsPrivate          bool      `json:"is_private" gorm:"default:true"`
	AdminLikedPoems    []Poem    `json:"admin_liked_poems" gorm:"many2many:admin_liked_poems;"`
	AdminBookmarkPoems []Poem    `json:"admin_bookmark_poems" gorm:"many2many:admin_bookmark_poems;"`
	UserBooksRead      []Book    `json:"user_books_read" gorm:"many2many:user_books_read;"`
	Comments           []Comment `json:"comments" gorm:"foreignKey:AdminID"`
}
