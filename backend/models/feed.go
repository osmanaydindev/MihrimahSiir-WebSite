package models

import "time"

// Akışın ihtiyaç duyduğu join tabloları.
//
// Beğeni/okundu ilişkileri bugüne kadar GORM'un otomatik ürettiği many2many
// tablolarıydı: ne zaman damgası ne de benzersizlik kısıtı vardı. Zaman
// damgası olmadan kronolojik akış kurulamıyor; kısıt olmadan da aynı şiiri
// iki kez beğenmek mükerrer satır ekleyip like_count'u şişiriyordu.
//
// Bu struct'lar db.SetupJoinTable ile mevcut tablolara bağlanıyor
// (bkz. database/connect.go) — tablo adları korunuyor, veri taşınmıyor.
//
// Eski satırlarda CreatedAt NULL kalır. Akış bunları göstermez: geçmişe ait
// bir beğeninin ne zaman yapıldığı bilinmiyor, uydurma bir tarihle akışa
// sokmak yanıltıcı olurdu.

// AdminLikedPoem, admin_liked_poems tablosunun açık karşılığı.
type AdminLikedPoem struct {
	AdminID   uint       `gorm:"primaryKey;autoIncrement:false"`
	PoemID    uint       `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt *time.Time `gorm:"index"`
}

func (AdminLikedPoem) TableName() string { return "admin_liked_poems" }

// AdminBookmarkPoem, admin_bookmark_poems tablosunun açık karşılığı.
type AdminBookmarkPoem struct {
	AdminID   uint       `gorm:"primaryKey;autoIncrement:false"`
	PoemID    uint       `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt *time.Time `gorm:"index"`
}

func (AdminBookmarkPoem) TableName() string { return "admin_bookmark_poems" }

// UserBookRead, user_books_read tablosunun açık karşılığı.
type UserBookRead struct {
	AdminID   uint       `gorm:"primaryKey;autoIncrement:false"`
	BookID    uint       `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt *time.Time `gorm:"index"`
}

func (UserBookRead) TableName() string { return "user_books_read" }

// ─── Alıntı beğenme / kaydetme ──────────────────────────────────────────────
// Şiirlerdeki beğeni/kaydetme yalnızca Poem için vardı; alıntılar (Comment)
// için hiçbir şey yoktu. Bileşik birincil anahtar, poem tablolarındaki
// mükerrer kayıt hatasının burada tekrarlanmasını engelliyor.

type CommentLike struct {
	AdminID   uint      `json:"admin_id"   gorm:"primaryKey;autoIncrement:false"`
	CommentID uint      `json:"comment_id" gorm:"primaryKey;autoIncrement:false;index"`
	CreatedAt time.Time `json:"created_at"`
}

func (CommentLike) TableName() string { return "comment_likes" }

type CommentSave struct {
	AdminID   uint      `json:"admin_id"   gorm:"primaryKey;autoIncrement:false"`
	CommentID uint      `json:"comment_id" gorm:"primaryKey;autoIncrement:false;index"`
	CreatedAt time.Time `json:"created_at"`
}

func (CommentSave) TableName() string { return "comment_saves" }
