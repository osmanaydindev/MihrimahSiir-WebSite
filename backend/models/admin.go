package models

import "time"

type Admin struct {
	ID                            uint       `json:"id" autoIncrement:"true"`
	Username                      string     `json:"username" gorm:"unique"`
	Email                         string     `json:"email" gorm:"unique"`
	Phone                         string     `json:"phone"`
	Password                      []byte     `json:"password" readOnly:"false"`
	RoleID                        uint       `json:"role_id"`
	ProfileImage                  string     `json:"profile_image"`
	IsPrivate                     bool       `json:"is_private" gorm:"default:true"`
	EmailVerified                 bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt               *time.Time `json:"email_verified_at"`
	EmailVerificationTokenHash    string     `json:"-" gorm:"index"`
	EmailVerificationTokenExpires *time.Time `json:"-"`
	PasswordResetTokenHash        string     `json:"-" gorm:"index"`
	PasswordResetTokenExpires     *time.Time `json:"-"`
	AdminLikedPoems               []Poem     `json:"admin_liked_poems" gorm:"many2many:admin_liked_poems;"`
	AdminBookmarkPoems            []Poem     `json:"admin_bookmark_poems" gorm:"many2many:admin_bookmark_poems;"`
	UserBooksRead                 []Book     `json:"user_books_read" gorm:"many2many:user_books_read;"`
	Comments                      []Comment  `json:"comments" gorm:"foreignKey:AdminID"`
}
