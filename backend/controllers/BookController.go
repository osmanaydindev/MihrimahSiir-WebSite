package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"backend/security"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// book's community: 1-özel (rol 1,2), 2-herkese açık, 3-sadece seçili kullanıcılar
func applyCommunityFilterForBook(db *gorm.DB, roleID uint, userID uint) *gorm.DB {
	return helpers.ApplyBookCommunityFilter(db, roleID, userID)
}

func CreateBook(c *fiber.Ctx) error {
	var payload bookWriteRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Geçersiz istek gövdesi",
		})
	}

	book := models.Book{IsDeleted: false}
	if err := applyBookPayload(&book, payload, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	book.CreatedAt = time.Now().Format("02-01-2006")
	book.Slug = helpers.UniqueBookSlug(book.Name, 0)

	if err := database.DB.Create(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Kitap oluşturulamadı",
		})
	}

	if err := helpers.SetBookVisibility(database.DB, book.ID, book.Community, payload.VisibleUserIDs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Kitap görünürlüğü kaydedilemedi",
		})
	}

	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getBooks(roleID, userID))
}
func GetBooks(c *fiber.Ctx) error {
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	books := getBooks(roleID, userID)
	return c.JSON(books)
}

// GetBooksPaginated returns paginated list of books with community and friendship filtering
func GetBooksPaginated(c *fiber.Ctx) error {
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	// Get pagination parameters and search query
	params := helpers.GetPaginationParams(c)
	search := c.Query("search", "")

	// Build base query
	baseQuery := database.DB.Model(&models.Book{}).Where("is_deleted = ?", false)
	baseQuery = applyCommunityFilterForBook(baseQuery, roleID, userID)

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR author ILIKE ?", searchPattern, searchPattern)
	}

	// Get total count
	var total int64
	baseQuery.Count(&total)

	// Get paginated books with community filtering
	var books []models.Book
	query := database.DB.Where("is_deleted = ?", false)
	query = applyCommunityFilterForBook(query, roleID, userID)

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR author ILIKE ?", searchPattern, searchPattern)
	}

	// Preload comments with Admin data
	query = query.Preload("Comments", "is_deleted = ?", false).
		Preload("Comments.Admin").
		Preload("AuthorData")

	query.Offset(params.Offset).
		Limit(params.Limit).
		Order("created_at DESC").
		Find(&books)

	// Filter comments by friendship for each book
	for i := range books {
		books[i].Comments = filterCommentsByFriendship(books[i].Comments, userID, roleID)
		// Ensure Comments is never nil
		if books[i].Comments == nil {
			books[i].Comments = []models.Comment{}
		}
	}

	// Create paginated response
	response := helpers.CreatePaginationResponse(books, total, params.Offset, params.Limit)

	fmt.Printf("[GetBooksPaginated] roleID=%d, userID=%d, search=%s, total=%d, returned=%d\n", roleID, userID, search, total, len(books))

	return c.JSON(response)
}

func GetBook(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	var book models.Book
	query := database.DB.Where("slug = ?", slug)
	query = applyCommunityFilterForBook(query, roleID, userID)
	query.Preload("AuthorData").
		Preload("Comments", "is_deleted = ?", false).
		Preload("Comments.Admin").
		First(&book)

	// Debug: Log before filtering
	fmt.Printf("[DEBUG GetBook] userID=%d, roleID=%d, slug=%s, comments_before=%d\n", userID, roleID, slug, len(book.Comments))
	for i, comment := range book.Comments {
		fmt.Printf("[DEBUG GetBook] Comment #%d: ID=%d, AdminID=%d, Admin=%s\n", i, comment.ID, comment.AdminID, comment.Admin.Username)
	}

	// Filter comments by friendship
	book.Comments = filterCommentsByFriendship(book.Comments, userID, roleID)

	// Ensure Comments is never nil (should be empty array instead)
	if book.Comments == nil {
		book.Comments = []models.Comment{}
	}

	// Debug: Log after filtering
	fmt.Printf("[DEBUG GetBook] comments_after=%d\n", len(book.Comments))
	for i, comment := range book.Comments {
		fmt.Printf("[DEBUG GetBook] Filtered Comment #%d: ID=%d, AdminID=%d, Admin=%s\n", i, comment.ID, comment.AdminID, comment.Admin.Username)
	}

	return c.JSON(book)
}

func getBooks(roleID uint, userID uint) []models.Book {
	var books []models.Book
	query := database.DB.Where("is_deleted = ?", false)
	query = applyCommunityFilterForBook(query, roleID, userID)
	query.Preload("AuthorData").
		Preload("Comments", "is_deleted = ?", false).
		Preload("Comments.Admin").
		Find(&books)

	// Filter comments by friendship for each book
	for i := range books {
		books[i].Comments = filterCommentsByFriendship(books[i].Comments, userID, roleID)
	}

	return books
}

func GetBookById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	var book models.Book
	query := database.DB.Table("books").Where("books.id = ?", id)
	query = applyCommunityFilterForBook(query, roleID, userID)
	result := query.Preload("AuthorData").First(&book)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	// Admin düzenleme ekranı seçili kullanıcı listesini de görebilmeli
	return c.JSON(fiber.Map{
		"id":               book.ID,
		"name":             book.Name,
		"author":           book.Author,
		"author_id":        book.AuthorID,
		"slug":             book.Slug,
		"image":            book.Image,
		"page":             book.Page,
		"isbn":             book.ISBN,
		"description":      book.Description,
		"community":        book.Community,
		"created_at":       book.CreatedAt,
		"author_data":      book.AuthorData,
		"visible_user_ids": helpers.GetBookVisibilityUserIDs(book.ID),
	})
}

func UpdateBook(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var book models.Book
	if err := database.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	var payload bookWriteRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	previousName := book.Name
	if err := applyBookPayload(&book, payload, false); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Slug'ı sadece isim gerçekten değiştiyse yenile; aksi halde mevcut
	// bağlantılar (ve paylaşılmış URL'ler) boşuna kırılır.
	if book.Name != previousName {
		book.Slug = helpers.UniqueBookSlug(book.Name, book.ID)
	}

	if err := database.DB.Save(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}

	// Görünürlük listesi sadece istemci açıkça gönderdiyse ya da seviye
	// 3'ten çıkıldıysa dokunulur; eski istemciler listeyi silmesin.
	if payload.VisibleUserIDs != nil || book.Community != 3 {
		if err := helpers.SetBookVisibility(database.DB, book.ID, book.Community, payload.VisibleUserIDs); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Kitap görünürlüğü kaydedilemedi",
			})
		}
	}

	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getBooks(roleID, userID))
}

// bookWriteRequest, admin kitap oluşturma/güncelleme gövdesi.
// Doğrudan models.Book'a BodyParser yapmak mass assignment'a açıktı
// (is_deleted, slug, id gibi alanlar dışarıdan set edilebiliyordu).
type bookWriteRequest struct {
	Name           string `json:"name"`
	Author         string `json:"author"`
	AuthorID       *uint  `json:"author_id"`
	Image          string `json:"image"`
	Description    string `json:"description"`
	ISBN           string `json:"isbn"`
	Page           int    `json:"page"`
	Community      int    `json:"community"`
	VisibleUserIDs []uint `json:"visible_user_ids"`
}

// applyBookPayload, gövdeyi doğrular/sanitize eder ve book üzerine uygular.
// isCreate=false iken boş alanlar "değiştirme" anlamına gelir (mevcut
// kısmi güncelleme davranışı korunuyor).
func applyBookPayload(book *models.Book, payload bookWriteRequest, isCreate bool) error {
	sanitizer := security.NewSanitizer()
	validator := security.NewValidator()

	name := sanitizer.SanitizeString(payload.Name, 500)
	if isCreate || name != "" {
		if err := validator.ValidateString("name", name, 1, 500, true); err != nil {
			return err
		}
		if sanitizer.ContainsDangerousContent(name) {
			return fmt.Errorf("kitap adı geçersiz karakterler içeriyor")
		}
		book.Name = name
	}

	if author := sanitizer.SanitizeString(payload.Author, 255); author != "" {
		book.Author = author
	}
	if payload.AuthorID != nil {
		book.AuthorID = payload.AuthorID
	}
	if image := sanitizer.SanitizeString(payload.Image, 1000); image != "" {
		book.Image = image
	}
	if payload.Description != "" {
		book.Description = sanitizer.StripTags(sanitizer.SanitizeString(payload.Description, 5000))
	}
	if payload.ISBN != "" {
		normalized, err := validator.NormalizeISBN(payload.ISBN)
		if err != nil {
			return err
		}
		book.ISBN = normalized
	}
	if payload.Page > 0 {
		if err := validator.ValidateInteger("page", payload.Page, 1, 100000); err != nil {
			return err
		}
		book.Page = payload.Page
	}
	if payload.Community > 0 {
		if err := validator.ValidateBookCommunity(payload.Community); err != nil {
			return err
		}
		book.Community = payload.Community
	} else if isCreate {
		book.Community = 2
	}

	if book.Community == 3 && len(payload.VisibleUserIDs) == 0 {
		return fmt.Errorf("görünürlük 'seçili kullanıcılar' iken en az bir kullanıcı seçilmeli")
	}

	return nil
}

func DeleteBook(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var book models.Book
	database.DB.Table("books").Where("id", id).Find(&book)
	book.IsDeleted = true
	database.DB.Save(&book)

	// Delete all comments associated with this book
	database.DB.Model(&models.Comment{}).Where("book_id = ?", id).Update("is_deleted", true)

	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getBooks(roleID, userID))
}

// filterCommentsByFriendship filters comments based on friendship relationships
func filterCommentsByFriendship(comments []models.Comment, userID uint, roleID uint) []models.Comment {
	fmt.Printf("[DEBUG filterComments] userID=%d, roleID=%d, total_comments=%d\n", userID, roleID, len(comments))

	// Admin can see all comments
	if roleID == 1 {
		fmt.Println("[DEBUG filterComments] User is admin, returning all comments")
		return comments
	}

	// Get friend IDs (includes self)
	friendIDs := GetFriendIDs(userID)
	fmt.Printf("[DEBUG filterComments] friendIDs=%v\n", friendIDs)

	// Filter comments to only include those from friends
	filteredComments := []models.Comment{}
	for _, comment := range comments {
		isFriend := false
		for _, friendID := range friendIDs {
			if comment.AdminID == friendID {
				isFriend = true
				break
			}
		}
		fmt.Printf("[DEBUG filterComments] Comment ID=%d, AdminID=%d, isFriend=%v\n", comment.ID, comment.AdminID, isFriend)
		if isFriend {
			filteredComments = append(filteredComments, comment)
		}
	}

	fmt.Printf("[DEBUG filterComments] Filtered result: %d comments\n", len(filteredComments))
	return filteredComments
}
