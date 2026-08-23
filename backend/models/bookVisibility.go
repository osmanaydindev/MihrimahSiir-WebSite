package models

// BookVisibility, community=3 (sadece seçili kullanıcılar) olan kitaplara
// hangi kullanıcıların erişebileceğini tutar. Bileşik birincil anahtar
// aynı kitap-kullanıcı çiftinin iki kez eklenmesini engeller.
type BookVisibility struct {
	BookID  uint `json:"book_id"  gorm:"primaryKey;autoIncrement:false"`
	AdminID uint `json:"admin_id" gorm:"primaryKey;autoIncrement:false;index"`
}

func (BookVisibility) TableName() string {
	return "book_visibilities"
}
