package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AddBookToReads(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	if uint(id) != userID {
		return c.SendStatus(404)
	}
	var user models.Admin
	if database.DB.
		Preload("AdminLikedPoems").
		Preload("AdminBookmarkPoems").
		Preload("UserBooksRead").
		Where("id", id).
		First(&user); user.ID == 0 {
		return c.SendStatus(404)
	}

	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.SendStatus(404)
	}

	err1 := database.DB.Model(&user).Association("UserBooksRead").Append(&book)
	if err1 != nil {
		return c.SendStatus(404)
	}

	database.DB.Save(&user)

	return c.JSON(user)
}
func DeleteBookFromReads(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	if uint(id) != userID {
		return c.SendStatus(404)
	}
	var user models.Admin
	if database.DB.
		Preload("AdminLikedPoems").
		Preload("AdminBookmarkPoems").
		Preload("UserBooksRead").
		Where("id", id).
		First(&user); user.ID == 0 {
		return c.SendStatus(404)
	}
	var book models.Book
	if err := c.BodyParser(&book); err != nil {
		return c.SendStatus(404)
	}

	err1 := database.DB.Model(&user).Association("UserBooksRead").Delete(&book)
	if err1 != nil {
		return c.SendStatus(404)
	}

	database.DB.Save(&user)

	return c.JSON(user)
}

// GetReadBooksIds returns only IDs of books that user has read (lightweight endpoint)
func GetReadBooksIds(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	if uint(id) != userID {
		return c.SendStatus(fiber.StatusForbidden)
	}

	var bookIDs []uint
	database.DB.Table("user_books_read").
		Where("admin_id = ?", userID).
		Pluck("book_id", &bookIDs)

	return c.JSON(bookIDs)
}

// GetReadsBooksPaginated returns paginated list of books that user has read
func GetReadsBooksPaginated(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	// Verify user is requesting their own books
	if uint(id) != userID {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// Get pagination parameters and search query
	params := helpers.GetPaginationParams(c)
	search := c.Query("search", "")

	// Get all read book IDs for this user
	var allBookIDs []uint
	database.DB.Table("user_books_read").
		Where("admin_id = ?", userID).
		Pluck("book_id", &allBookIDs)

	// If no read books, return empty result
	if len(allBookIDs) == 0 {
		response := helpers.CreatePaginationResponse([]models.Book{}, 0, params.Offset, params.Limit)
		return c.JSON(response)
	}

	// Build query for books
	query := database.DB.Where("id IN ? AND is_deleted = ?", allBookIDs, false)

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR author ILIKE ?", searchPattern, searchPattern)
	}

	// Count total matching books
	var total int64
	query.Model(&models.Book{}).Count(&total)

	// Fetch paginated books
	books := []models.Book{}
	query = database.DB.Where("id IN ? AND is_deleted = ?", allBookIDs, false)

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR author ILIKE ?", searchPattern, searchPattern)
	}

	// Preload comments with Admin data and AuthorData
	query = query.Preload("Comments", "is_deleted = ?", false).
		Preload("Comments.Admin").
		Preload("AuthorData")

	query.Offset(params.Offset).
		Limit(params.Limit).
		Find(&books)

	// Filter comments by friendship for each book
	for i := range books {
		books[i].Comments = filterCommentsByFriendship(books[i].Comments, userID, roleID)
		if books[i].Comments == nil {
			books[i].Comments = []models.Comment{}
		}
	}

	// Create paginated response
	response := helpers.CreatePaginationResponse(books, total, params.Offset, params.Limit)

	fmt.Printf("[GetReadsBooksPaginated] userID=%d, search=%s, total=%d, returned=%d\n", userID, search, total, len(books))

	return c.JSON(response)
}

// GetNotReadsBooksPaginated returns paginated list of books that user wants to read (all books minus read books)
func GetNotReadsBooksPaginated(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	// Verify user is requesting their own books
	if uint(id) != userID {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// Get pagination parameters and search query
	params := helpers.GetPaginationParams(c)
	search := c.Query("search", "")

	// Get IDs of books user has already read
	var readBookIDs []uint
	database.DB.Table("user_books_read").
		Where("admin_id = ?", userID).
		Pluck("book_id", &readBookIDs)

	// Build base query for books user hasn't read yet
	baseQuery := database.DB.Model(&models.Book{}).Where("is_deleted = ?", false)
	baseQuery = applyCommunityFilterForBook(baseQuery, roleID)
	if len(readBookIDs) > 0 {
		baseQuery = baseQuery.Where("id NOT IN ?", readBookIDs)
	}

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		baseQuery = baseQuery.Where("name ILIKE ? OR author ILIKE ?", searchPattern, searchPattern)
	}

	// Count total unread books
	var total int64
	baseQuery.Count(&total)

	// Get paginated unread books
	books := []models.Book{}
	query := database.DB.Where("is_deleted = ?", false)
	query = applyCommunityFilterForBook(query, roleID)
	if len(readBookIDs) > 0 {
		query = query.Where("id NOT IN ?", readBookIDs)
	}

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("name ILIKE ? OR author ILIKE ?", searchPattern, searchPattern)
	}

	// Preload comments with Admin data and AuthorData
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
		if books[i].Comments == nil {
			books[i].Comments = []models.Comment{}
		}
	}

	// Create paginated response
	response := helpers.CreatePaginationResponse(books, total, params.Offset, params.Limit)

	fmt.Printf("[GetNotReadsBooksPaginated] userID=%d, search=%s, total=%d, returned=%d\n", userID, search, total, len(books))

	return c.JSON(response)
}
