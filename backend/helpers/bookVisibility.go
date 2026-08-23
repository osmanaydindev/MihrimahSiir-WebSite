package helpers

import (
	"backend/database"
	"backend/models"
	"gorm.io/gorm"
)

// ApplyBookCommunityFilter, kitap sorgularına görünürlük kısıtını ekler.
//
//	community 1 = özel      -> sadece role_id 1 ve 2
//	community 2 = herkese açık
//	community 3 = seçili    -> book_visibilities tablosundaki kullanıcılar
//
// Sorgu "books" tablosunu içermeli (doğrudan ya da JOIN ile); kolonlar
// bilerek "books." ile niteleniyor ki JOIN'li sorgularda da çalışsın.
// Dış parantezler zorunlu: bu ifade çağıran tarafta is_deleted / arama
// koşullarıyla AND'lendiği için OR'un kapsamı korunmalı.
func ApplyBookCommunityFilter(db *gorm.DB, roleID uint, userID uint) *gorm.DB {
	if roleID == 1 || roleID == 2 {
		return db
	}

	visibleBookIDs := database.DB.Model(&models.BookVisibility{}).
		Select("book_id").
		Where("admin_id = ?", userID)

	return db.Where(
		"(books.community = ? OR (books.community = ? AND books.id IN (?)))",
		2, 3, visibleBookIDs,
	)
}

// SetBookVisibility, bir kitabın seçili-kullanıcı listesini değiştirir.
// community 3 değilse liste tamamen temizlenir, böylece görünürlük
// seviyesi düşürülen kitaplarda öksüz satır kalmaz.
func SetBookVisibility(tx *gorm.DB, bookID uint, community int, userIDs []uint) error {
	if err := tx.Where("book_id = ?", bookID).Delete(&models.BookVisibility{}).Error; err != nil {
		return err
	}

	if community != 3 || len(userIDs) == 0 {
		return nil
	}

	rows := make([]models.BookVisibility, 0, len(userIDs))
	for _, id := range userIDs {
		rows = append(rows, models.BookVisibility{BookID: bookID, AdminID: id})
	}
	return tx.Create(&rows).Error
}

// GetBookVisibilityUserIDs, kitabın seçili kullanıcı id'lerini döner.
func GetBookVisibilityUserIDs(bookID uint) []uint {
	ids := []uint{}
	database.DB.Model(&models.BookVisibility{}).
		Where("book_id = ?", bookID).
		Pluck("admin_id", &ids)
	return ids
}
