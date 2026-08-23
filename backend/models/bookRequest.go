package models

import "time"

// Talep durumları. Friendship.Status ile aynı desen (varchar + default).
const (
	BookRequestPending  = "pending"
	BookRequestApproved = "approved"
	BookRequestRejected = "rejected"
)

// BookRequest, kullanıcının ISBN ile açtığı kitap ekleme talebi.
// Open Library verisi talep anında çekilip burada saklanır (approval'da
// yeniden çekilmez): admin listesi N talebi N dış çağrı olmadan
// gösterebilsin, onay dış servis çökse de çalışsın ve admin'in
// incelediği veriyle kaydedilen veri aynı olsun diye.
type BookRequest struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	UserID uint   `json:"user_id" gorm:"not null;index"`
	ISBN   string `json:"isbn" gorm:"type:varchar(13);not null;index"` // her zaman ISBN-13
	Status string `json:"status" gorm:"type:varchar(20);not null;default:'pending';index"`

	// Open Library anlık görüntüsü
	FetchedTitle       string `json:"fetched_title" gorm:"type:varchar(500)"`
	FetchedAuthors     string `json:"fetched_authors" gorm:"type:varchar(500)"` // virgülle ayrılmış
	FetchedPages       int    `json:"fetched_pages"`
	FetchedCoverURL    string `json:"fetched_cover_url" gorm:"type:varchar(1000)"`
	FetchedDescription string `json:"fetched_description" gorm:"type:text"`
	FetchedPublisher   string `json:"fetched_publisher" gorm:"type:varchar(255)"`
	FetchedPublishDate string `json:"fetched_publish_date" gorm:"type:varchar(64)"`
	OpenLibraryKey     string `json:"open_library_key" gorm:"type:varchar(64)"`
	MetadataFound      bool   `json:"metadata_found" gorm:"default:false"`

	UserNote      string     `json:"user_note" gorm:"type:varchar(500)"`
	AdminNote     string     `json:"admin_note" gorm:"type:varchar(500)"` // red gerekçesi
	ReviewedBy    *uint      `json:"reviewed_by"`
	ReviewedAt    *time.Time `json:"reviewed_at"`
	CreatedBookID *uint      `json:"created_book_id" gorm:"index"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User        Admin `json:"user" gorm:"foreignKey:UserID"`
	CreatedBook *Book `json:"created_book,omitempty" gorm:"foreignKey:CreatedBookID"`
}

func (BookRequest) TableName() string {
	return "book_requests"
}
