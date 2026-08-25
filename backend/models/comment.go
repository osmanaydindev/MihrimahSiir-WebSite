package models

import "time"

type Comment struct {
	ID      uint  `json:"id" autoIncrement:"true"`
	AdminID uint  `json:"admin_id"`
	Admin   Admin `json:"admin" gorm:"foreignKey:AdminID"`
	BookID  uint  `json:"book_id"`
	Book    Book  `json:"book" gorm:"foreignKey:BookID"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Page    *int   `json:"page"` // Optional page number
	IsDeleted bool `json:"is_deleted" gorm:"default:false"`

	// Akış kronolojik sıralama için bunu kullanıyor. Model bugüne kadar
	// zaman damgası taşımıyordu; listeler "id DESC" ile sıralanıyordu ve
	// commentList.vue'daki formatDate hep boş basıyordu.
	// Mevcut satırlar EnsureIndexes() içinde tek seferlik dolduruluyor.
	CreatedAt time.Time `json:"created_at"`

	// Bu yorumun, "public profil" kuralının yürürlüğe girmesinden SONRA
	// yazılıp yazılmadığı.
	//
	// Bu alandan önce yorumlar yalnızca arkadaşlara görünüyordu. Kuralı
	// geçmişe dönük genişletmek, kullanıcıların "sadece arkadaşlarım görecek"
	// varsayımıyla yazdığı notları bir anda herkese açardı. Eski satırlar
	// false kalır ve asla public profil üzerinden görünmez; AddComment
	// bundan sonra true yazar.
	AllowPublicFeed bool `json:"-" gorm:"default:false"`
}
