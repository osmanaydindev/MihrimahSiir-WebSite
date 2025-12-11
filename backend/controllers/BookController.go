package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// book's community: 1-private, 2-public
func applyCommunityFilterForBook(db *gorm.DB, roleID uint) *gorm.DB {
	if roleID == 1 || roleID == 2 {
		return db
	}
	return db.Where("community = ?", 2)
}

func CreateBook(c *fiber.Ctx) error {
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return err
	}
	book.IsDeleted = false
	x := map[rune]rune{
		' ':  '-',
		'ç':  'c',
		'ğ':  'g',
		'ı':  'i',
		'ö':  'o',
		'ş':  's',
		'ü':  'u',
		'Ç':  'C',
		'Ğ':  'G',
		'İ':  'I',
		'Ö':  'O',
		'Ş':  'S',
		'Ü':  'U',
		'â':  'a',
		'Â':  'A',
		'\'': '-',
	}
	slug := replaceChars(book.Name, x)
	book.CreatedAt = time.Now().Format("02-01-2006")
	book.Slug = strings.ToLower(slug)
	database.DB.Create(&book)

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
	baseQuery = applyCommunityFilterForBook(baseQuery, roleID)

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
	query = applyCommunityFilterForBook(query, roleID)

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
	query = applyCommunityFilterForBook(query, roleID)
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
	query = applyCommunityFilterForBook(query, roleID)
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
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	var book models.Book
	query := database.DB.Table("books").Where("id", id)
	query = applyCommunityFilterForBook(query, roleID)
	result := query.Preload("AuthorData").First(&book)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	return c.JSON(book)
}

func UpdateBook(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var book models.Book
	if err := database.DB.Where("id = ?", id).First(&book).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Book not found",
		})
	}

	var updateData models.Book
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if updateData.Name != "" {
		book.Name = updateData.Name
		// Regenerate slug if name changed
		x := map[rune]rune{
			' ':  '-',
			'ç':  'c',
			'ğ':  'g',
			'ı':  'i',
			'ö':  'o',
			'ş':  's',
			'ü':  'u',
			'Ç':  'C',
			'Ğ':  'G',
			'İ':  'I',
			'Ö':  'O',
			'Ş':  'S',
			'Ü':  'U',
			'â':  'a',
			'Â':  'A',
			'\'': '-',
		}
		slug := replaceChars(book.Name, x)
		book.Slug = strings.ToLower(slug)
	}
	if updateData.Author != "" {
		book.Author = updateData.Author
	}
	if updateData.AuthorID != nil {
		book.AuthorID = updateData.AuthorID
	}
	if updateData.Image != "" {
		book.Image = updateData.Image
	}
	if updateData.Page > 0 {
		book.Page = updateData.Page
	}
	if updateData.Community > 0 {
		book.Community = updateData.Community
	}

	if err := database.DB.Save(&book).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update book",
		})
	}

	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getBooks(roleID, userID))
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
