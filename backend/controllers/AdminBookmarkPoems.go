package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AddBookmark(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var user models.Admin
	if database.DB.
		Preload("AdminLikedPoems").
		Preload("AdminBookmarkPoems").
		Preload("UserBooksRead").
		Where("id", id).
		First(&user); user.ID == 0 {
		return c.SendStatus(404)
	}
	userID := GetUserId(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}

	var poem models.Poem
	if err := c.BodyParser(&poem); err != nil {
		return c.SendStatus(404)
	}

	err1 := database.DB.Model(&user).Association("AdminBookmarkPoems").Append(&poem)
	if err1 != nil {
		return c.SendStatus(404)
	}

	database.DB.Save(&user)

	return c.JSON(user)
}

func UndoBookmark(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var user models.Admin
	if database.DB.
		Preload("AdminLikedPoems").
		Preload("AdminBookmarkPoems").
		Preload("UserBooksRead").
		Where("id", id).
		First(&user); user.ID == 0 {
		return c.SendStatus(404)
	}
	userID := GetUserId(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	var poem models.Poem
	if err := c.BodyParser(&poem); err != nil {
		return c.SendStatus(404)
	}

	err1 := database.DB.Model(&user).Association("AdminBookmarkPoems").Delete(&poem)
	if err1 != nil {
		return c.SendStatus(404)
	}

	database.DB.Save(&user)

	return c.JSON(user)
}

// GetBookmarksIdByAdminId returns only IDs of poems that user bookmarked (optimized)
func GetBookmarksIdByAdminId(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}

	var poemIDs []uint
	database.DB.Table("admin_bookmark_poems").
		Where("admin_id = ?", userID).
		Pluck("poem_id", &poemIDs)

	return c.JSON(poemIDs)
}

// GetBookmarksPaginated returns paginated list of poems that user bookmarked
func GetBookmarksPaginated(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)

	if uint(id) != userID {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// Get pagination parameters and search query
	params := helpers.GetPaginationParams(c)
	search := c.Query("search", "")

	// Get all bookmarked poem IDs for this user
	var allPoemIDs []uint
	database.DB.Table("admin_bookmark_poems").
		Where("admin_id = ?", userID).
		Pluck("poem_id", &allPoemIDs)

	// If no bookmarked poems, return empty result
	if len(allPoemIDs) == 0 {
		response := helpers.CreatePaginationResponse([]models.Poem{}, 0, params.Offset, params.Limit)
		return c.JSON(response)
	}

	// Build query for poems
	query := database.DB.Where("id IN ? AND is_deleted = ?", allPoemIDs, false)

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchPattern, searchPattern)
	}

	// Count total matching poems
	var total int64
	query.Model(&models.Poem{}).Count(&total)

	// Fetch paginated poems
	poems := []models.Poem{}
	query.Offset(params.Offset).
		Limit(params.Limit).
		Find(&poems)

	// Create paginated response
	response := helpers.CreatePaginationResponse(poems, total, params.Offset, params.Limit)

	fmt.Printf("[GetBookmarksPaginated] userID=%d, search=%s, total=%d, returned=%d\n", userID, search, total, len(poems))

	return c.JSON(response)
}
