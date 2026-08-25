package helpers

import (
	"backend/database"
	"backend/models"
)

// AreFriends, iki kullanıcının kabul edilmiş arkadaşlığı olup olmadığını söyler.
// Kişi kendisiyle her zaman "arkadaş" sayılır (kendi içeriğini görebilsin).
//
// Not: controllers.AreFriends ile aynı işi yapar. Buraya kopyalandı çünkü
// helpers paketi controllers'ı import edemez (import döngüsü). Uzun vadede
// controllers tarafındaki sürüm bunu çağırmalı.
func AreFriends(userID1, userID2 uint) bool {
	if userID1 == userID2 {
		return true
	}

	var friendship models.Friendship
	err := database.DB.Where(
		"((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?",
		userID1, userID2, userID2, userID1, "accepted",
	).First(&friendship).Error

	return err == nil
}

// FriendIDs, kullanıcının arkadaş id'lerini döner — kendisi de listeye dahil.
func FriendIDs(userID uint) []uint {
	var friendships []models.Friendship
	database.DB.Where(
		"(user_id = ? OR friend_id = ?) AND status = ?",
		userID, userID, "accepted",
	).Find(&friendships)

	ids := []uint{userID}
	for _, fs := range friendships {
		if fs.UserID == userID {
			ids = append(ids, fs.FriendID)
		} else {
			ids = append(ids, fs.UserID)
		}
	}
	return ids
}

// PublicUserIDs, profili herkese açık olan kullanıcıların id'lerini döner.
func PublicUserIDs() []uint {
	ids := []uint{}
	database.DB.Model(&models.Admin{}).
		Where("is_private = ?", false).
		Pluck("id", &ids)
	return ids
}

// CanViewUserContent, izleyicinin hedef kullanıcının içeriğini görüp
// göremeyeceğini söyler.
//
// Bu ifade AdminController içinde üç ayrı yerde birebir kopyalanmıştı
// (GetUserProfile, GetUserLikedPoems, GetUserReadBooks); tek yere alındı.
func CanViewUserContent(viewerID uint, target models.Admin) bool {
	if viewerID == target.ID {
		return true
	}
	if !target.IsPrivate {
		return true
	}
	return AreFriends(viewerID, target.ID)
}
