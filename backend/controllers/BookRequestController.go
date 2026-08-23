package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"backend/security"
	"backend/services/mail"
	"backend/services/openlibrary"
	ws "backend/websocket"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Kötüye kullanım tavanları. Rate limiter bellekte olduğu için
// restart'ta sıfırlanıyor; asıl koruma bu DB tabanlı sayımlar.
const (
	maxPendingRequestsPerUser = 3
	maxRequestsPerDay         = 10
	rejectionCooldownDays     = 30
)

type createBookRequestPayload struct {
	ISBN string `json:"isbn"`
	Note string `json:"note"`
}

// CreateBookRequest, kullanıcının ISBN ile açtığı talebi kaydeder.
// Kontrol sırası bilinçli: ucuz olanlar (checksum, indeksli SELECT'ler)
// önce çalışır, dış API çağrısı ve mail ancak hepsi geçilirse tetiklenir.
func CreateBookRequest(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Oturum bulunamadı",
		})
	}

	var payload createBookRequestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Geçersiz istek gövdesi",
		})
	}

	validator := security.NewValidator()
	sanitizer := security.NewSanitizer()

	// L2a — checksum. Rastgele hane dizisi burada elenir, soket açılmaz.
	isbn, err := validator.NormalizeISBN(payload.ISBN)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Geçersiz ISBN. Lütfen kitabın arkasındaki 13 haneli numarayı girin.",
		})
	}

	note := sanitizer.SanitizeString(payload.Note, 500)
	if sanitizer.ContainsDangerousContent(note) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Not alanı geçersiz içerik barındırıyor.",
		})
	}

	if status, message := checkRequestLimits(userID, isbn); message != "" {
		return c.Status(status).JSON(fiber.Map{"message": message})
	}

	// L4 — dış çağrı. Başarısız olursa talep yine oluşur; admin
	// bilgileri panelden elle girebilir.
	request := models.BookRequest{
		UserID:   userID,
		ISBN:     isbn,
		Status:   models.BookRequestPending,
		UserNote: note,
	}

	meta, fetchErr := openlibrary.Default().FetchByISBN(c.UserContext(),isbn)
	if fetchErr != nil {
		if !errors.Is(fetchErr, openlibrary.ErrNotFound) {
			log.Printf("[book-request] Open Library çağrısı başarısız (%s): %v", isbn, fetchErr)
		}
	} else {
		request.MetadataFound = true
		request.FetchedTitle = sanitizer.SanitizeString(meta.Title, 500)
		request.FetchedAuthors = sanitizer.SanitizeString(strings.Join(meta.Authors, ", "), 500)
		request.FetchedPages = meta.NumberOfPages
		request.FetchedCoverURL = meta.CoverURL
		request.FetchedDescription = sanitizer.SanitizeString(meta.Description, 5000)
		request.FetchedPublisher = sanitizer.SanitizeString(meta.Publisher, 255)
		request.FetchedPublishDate = sanitizer.SanitizeString(meta.PublishDate, 64)
		request.OpenLibraryKey = meta.EditionKey
	}

	if err := database.DB.Create(&request).Error; err != nil {
		// Partial unique index yarışı: iki kullanıcı aynı anda aynı ISBN'i istedi.
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Bu kitap için zaten bekleyen bir istek var.",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Kitap isteği oluşturulamadı",
		})
	}

	// L5 — mail. Hata isteği asla düşürmez.
	var requester models.Admin
	database.DB.Select("username, email").Where("id = ?", userID).First(&requester)
	mail.NotifyAdminNewBookRequest(mail.BookRequestInfo{
		Username:      requester.Username,
		Title:         request.FetchedTitle,
		Authors:       request.FetchedAuthors,
		ISBN:          request.ISBN,
		Publisher:     request.FetchedPublisher,
		PublishDate:   request.FetchedPublishDate,
		Pages:         request.FetchedPages,
		CoverURL:      request.FetchedCoverURL,
		UserNote:      request.UserNote,
		MetadataFound: request.MetadataFound,
	})

	message := "Kitap isteğin alındı, en kısa sürede değerlendirilecek."
	if !request.MetadataFound {
		message = "Kitap isteğin alındı. Bu ISBN için otomatik bilgi bulunamadı, bilgileri yönetici elle girecek."
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": message,
		"request": request,
	})
}

// checkRequestLimits, L2b–L3b katmanlarını uygular.
// Boş mesaj = kontrolden geçti.
func checkRequestLimits(userID uint, isbn string) (int, string) {
	// L2b — kitap zaten var mı
	var existingBooks int64
	database.DB.Model(&models.Book{}).
		Where("isbn = ? AND is_deleted = ?", isbn, false).
		Count(&existingBooks)
	if existingBooks > 0 {
		return fiber.StatusConflict, "Bu kitap zaten kütüphanede mevcut."
	}

	// L2c — aynı ISBN için bekleyen talep
	var pending models.BookRequest
	err := database.DB.Where("isbn = ? AND status = ?", isbn, models.BookRequestPending).First(&pending).Error
	if err == nil {
		if pending.UserID == userID {
			return fiber.StatusConflict, "Bu kitabı zaten istediniz, isteğiniz inceleniyor."
		}
		return fiber.StatusConflict, "Bu kitap için zaten bekleyen bir istek var."
	}

	// L2d — reddedilmiş talep için bekleme süresi. Bu olmadan kullanıcı
	// reddet -> tekrar iste döngüsüyle mail kotasını tüketebilir.
	var rejected models.BookRequest
	cooldownStart := time.Now().AddDate(0, 0, -rejectionCooldownDays)
	err = database.DB.
		Where("isbn = ? AND user_id = ? AND status = ? AND updated_at >= ?",
			isbn, userID, models.BookRequestRejected, cooldownStart).
		Order("updated_at DESC").
		First(&rejected).Error
	if err == nil {
		message := "Bu kitap daha önce reddedildi."
		if rejected.AdminNote != "" {
			message = fmt.Sprintf("Bu kitap daha önce reddedildi. Gerekçe: %s", rejected.AdminNote)
		}
		return fiber.StatusConflict, message
	}

	// L3a — eşzamanlı bekleyen talep sayısı
	var pendingCount int64
	database.DB.Model(&models.BookRequest{}).
		Where("user_id = ? AND status = ?", userID, models.BookRequestPending).
		Count(&pendingCount)
	if pendingCount >= maxPendingRequestsPerUser {
		return fiber.StatusTooManyRequests,
			fmt.Sprintf("Aynı anda en fazla %d bekleyen kitap isteğiniz olabilir.", maxPendingRequestsPerUser)
	}

	// L3b — günlük talep sayısı
	var dailyCount int64
	database.DB.Model(&models.BookRequest{}).
		Where("user_id = ? AND created_at >= ?", userID, time.Now().Add(-24*time.Hour)).
		Count(&dailyCount)
	if dailyCount >= maxRequestsPerDay {
		return fiber.StatusTooManyRequests,
			"Günlük kitap isteği limitine ulaştınız. Lütfen yarın tekrar deneyin."
	}

	return 0, ""
}

// GetMyBookRequests, kullanıcının kendi taleplerini döner.
func GetMyBookRequests(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Oturum bulunamadı",
		})
	}

	requests := []models.BookRequest{}
	database.DB.
		Preload("CreatedBook").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(100).
		Find(&requests)

	return c.JSON(requests)
}

// CancelBookRequest, kullanıcının kendi bekleyen talebini siler.
func CancelBookRequest(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Oturum bulunamadı",
		})
	}
	id, _ := strconv.Atoi(c.Params("id"))

	var request models.BookRequest
	if err := database.DB.First(&request, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "İstek bulunamadı",
		})
	}
	if request.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Sadece kendi isteğinizi iptal edebilirsiniz",
		})
	}
	if request.Status != models.BookRequestPending {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Bu istek zaten sonuçlandırılmış.",
		})
	}

	if err := database.DB.Delete(&request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "İstek iptal edilemedi",
		})
	}

	return c.JSON(fiber.Map{"message": "İstek iptal edildi"})
}

// GetBookRequests, admin listesi (sayfalı, duruma göre filtreli).
func GetBookRequests(c *fiber.Ctx) error {
	params := helpers.GetPaginationParams(c)
	status := c.Query("status", "")

	query := database.DB.Model(&models.BookRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	requests := []models.BookRequest{}
	dataQuery := database.DB.Preload("User").Preload("CreatedBook")
	if status != "" {
		dataQuery = dataQuery.Where("status = ?", status)
	}
	dataQuery.
		Offset(params.Offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&requests)

	// Parola hash'i listede dolaşmasın
	for i := range requests {
		requests[i].User.Password = nil
	}

	return c.JSON(helpers.CreatePaginationResponse(requests, total, params.Offset, params.Limit))
}

// GetBookRequestCount, panel rozeti için bekleyen talep sayısı.
func GetBookRequestCount(c *fiber.Ctx) error {
	var pending int64
	database.DB.Model(&models.BookRequest{}).
		Where("status = ?", models.BookRequestPending).
		Count(&pending)
	return c.JSON(fiber.Map{"pending": pending})
}

type approveBookRequestPayload struct {
	Name           string `json:"name"`
	Author         string `json:"author"`
	AuthorID       *uint  `json:"author_id"`
	Image          string `json:"image"`
	Description    string `json:"description"`
	Page           int    `json:"page"`
	Community      int    `json:"community"`
	VisibleUserIDs []uint `json:"visible_user_ids"`
}

// ApproveBookRequest, talebi onaylayıp kitabı oluşturur.
// Kitap oluşturma, görünürlük ve talebin durumu tek transaction'da.
func ApproveBookRequest(c *fiber.Ctx) error {
	adminID := GetUserId(c)
	id, _ := strconv.Atoi(c.Params("id"))

	var payload approveBookRequestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Geçersiz istek gövdesi",
		})
	}

	var request models.BookRequest
	if err := database.DB.Preload("User").First(&request, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "İstek bulunamadı",
		})
	}
	if request.Status != models.BookRequestPending {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Bu istek zaten sonuçlandırılmış.",
		})
	}

	validator := security.NewValidator()
	sanitizer := security.NewSanitizer()

	name := sanitizer.SanitizeString(payload.Name, 500)
	if err := validator.ValidateString("name", name, 1, 500, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Kitap adı zorunlu.",
		})
	}
	if sanitizer.ContainsDangerousContent(name) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Kitap adı geçersiz karakterler içeriyor.",
		})
	}

	community := payload.Community
	if community == 0 {
		community = 2
	}
	if err := validator.ValidateBookCommunity(community); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Geçersiz görünürlük değeri.",
		})
	}

	visibleUserIDs, err := resolveVisibleUserIDs(community, payload.VisibleUserIDs, request.UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	book := models.Book{
		Name:        name,
		Author:      sanitizer.SanitizeString(payload.Author, 255),
		AuthorID:    payload.AuthorID,
		Image:       sanitizer.SanitizeString(payload.Image, 1000),
		Description: sanitizer.StripTags(sanitizer.SanitizeString(payload.Description, 5000)),
		ISBN:        request.ISBN,
		Page:        payload.Page,
		Community:   community,
		IsDeleted:   false,
		CreatedAt:   time.Now().Format("02-01-2006"),
		Slug:        helpers.UniqueBookSlug(name, 0),
	}

	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&book).Error; err != nil {
			return err
		}
		if err := helpers.SetBookVisibility(tx, book.ID, community, visibleUserIDs); err != nil {
			return err
		}

		now := time.Now()
		request.Status = models.BookRequestApproved
		request.CreatedBookID = &book.ID
		request.ReviewedBy = &adminID
		request.ReviewedAt = &now
		return tx.Save(&request).Error
	})
	if txErr != nil {
		log.Printf("[book-request] onay başarısız (id=%d): %v", id, txErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Kitap oluşturulamadı",
		})
	}

	// Bildirimler transaction dışında: mail/websocket hatası onayı düşürmez.
	mail.NotifyUserBookRequestApproved(mail.BookRequestInfo{
		Username:    request.User.Username,
		UserEmail:   request.User.Email,
		Title:       book.Name,
		Authors:     request.FetchedAuthors,
		ISBN:        book.ISBN,
		Publisher:   request.FetchedPublisher,
		PublishDate: request.FetchedPublishDate,
		Pages:       book.Page,
		CoverURL:    request.FetchedCoverURL,
		BookSlug:    book.Slug,
	})
	ws.GlobalHub.SendToUser(request.UserID, "book_request_reviewed", map[string]interface{}{
		"status": models.BookRequestApproved,
		"title":  book.Name,
		"slug":   book.Slug,
	})

	return c.JSON(fiber.Map{
		"message": "Kitap isteği onaylandı ve kitap oluşturuldu.",
		"book":    book,
	})
}

// resolveVisibleUserIDs, seçili kullanıcı listesini doğrular:
// tekilleştirir, varlıklarını kontrol eder ve talep sahibini ekler
// (yoksa kullanıcı kendi istediği kitabı göremez).
func resolveVisibleUserIDs(community int, requested []uint, requesterID uint) ([]uint, error) {
	if community != 3 {
		return nil, nil
	}
	if len(requested) == 0 {
		return nil, errors.New("Görünürlük 'seçili kullanıcılar' iken en az bir kullanıcı seçilmeli.")
	}
	if len(requested) > 200 {
		return nil, errors.New("En fazla 200 kullanıcı seçebilirsiniz.")
	}

	unique := make(map[uint]struct{}, len(requested)+1)
	for _, id := range requested {
		if id > 0 {
			unique[id] = struct{}{}
		}
	}
	unique[requesterID] = struct{}{}

	ids := make([]uint, 0, len(unique))
	for id := range unique {
		ids = append(ids, id)
	}

	var existing int64
	database.DB.Model(&models.Admin{}).Where("id IN ?", ids).Count(&existing)
	if int(existing) != len(ids) {
		return nil, errors.New("Seçilen kullanıcılardan bazıları bulunamadı.")
	}

	return ids, nil
}

type rejectBookRequestPayload struct {
	AdminNote string `json:"admin_note"`
}

// RejectBookRequest, talebi reddeder ve gerekçeyi saklar.
func RejectBookRequest(c *fiber.Ctx) error {
	adminID := GetUserId(c)
	id, _ := strconv.Atoi(c.Params("id"))

	var payload rejectBookRequestPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Geçersiz istek gövdesi",
		})
	}

	var request models.BookRequest
	if err := database.DB.Preload("User").First(&request, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "İstek bulunamadı",
		})
	}
	if request.Status != models.BookRequestPending {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"message": "Bu istek zaten sonuçlandırılmış.",
		})
	}

	sanitizer := security.NewSanitizer()
	note := sanitizer.StripTags(sanitizer.SanitizeString(payload.AdminNote, 500))

	now := time.Now()
	request.Status = models.BookRequestRejected
	request.AdminNote = note
	request.ReviewedBy = &adminID
	request.ReviewedAt = &now

	if err := database.DB.Save(&request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "İstek reddedilemedi",
		})
	}

	mail.NotifyUserBookRequestRejected(mail.BookRequestInfo{
		Username:  request.User.Username,
		UserEmail: request.User.Email,
		Title:     request.FetchedTitle,
		Authors:   request.FetchedAuthors,
		ISBN:      request.ISBN,
		AdminNote: note,
	})
	ws.GlobalHub.SendToUser(request.UserID, "book_request_reviewed", map[string]interface{}{
		"status": models.BookRequestRejected,
		"title":  request.FetchedTitle,
	})

	return c.JSON(fiber.Map{"message": "Kitap isteği reddedildi."})
}

// RefreshBookRequest, Open Library anlık görüntüsünü yeniden çeker.
// Admin-only olduğu için saldırı yüzeyi değil.
func RefreshBookRequest(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var request models.BookRequest
	if err := database.DB.First(&request, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "İstek bulunamadı",
		})
	}

	meta, err := openlibrary.Default().FetchByISBN(c.UserContext(),request.ISBN)
	if err != nil {
		if errors.Is(err, openlibrary.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Open Library'de bu ISBN için kayıt bulunamadı.",
			})
		}
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "Open Library'ye şu anda ulaşılamıyor.",
		})
	}

	sanitizer := security.NewSanitizer()
	request.MetadataFound = true
	request.FetchedTitle = sanitizer.SanitizeString(meta.Title, 500)
	request.FetchedAuthors = sanitizer.SanitizeString(strings.Join(meta.Authors, ", "), 500)
	request.FetchedPages = meta.NumberOfPages
	request.FetchedCoverURL = meta.CoverURL
	request.FetchedDescription = sanitizer.SanitizeString(meta.Description, 5000)
	request.FetchedPublisher = sanitizer.SanitizeString(meta.Publisher, 255)
	request.FetchedPublishDate = sanitizer.SanitizeString(meta.PublishDate, 64)
	request.OpenLibraryKey = meta.EditionKey

	if err := database.DB.Save(&request).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "İstek güncellenemedi",
		})
	}

	return c.JSON(fiber.Map{"request": request})
}
