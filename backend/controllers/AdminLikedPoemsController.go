package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AddLikedPoem(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	var user models.Admin
	if database.DB.Preload("AdminLikedPoems").Preload("AdminBookmarkPoems").Where("id", id).First(&user); user.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"error": "Admin not found with id: " + c.Params("id"),
		})
	}
	var poem models.Poem
	if err := c.BodyParser(&poem); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid poem data: " + err.Error(),
		})
	}

	err1 := database.DB.Model(&user).Association("AdminLikedPoems").Append(&poem)
	if err1 != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to add liked poem: " + err1.Error(),
		})
	}

	database.DB.Save(&user)

	return c.JSON(user)
}
func UndoLikedPoem(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	var user models.Admin
	if database.DB.Preload("AdminLikedPoems").Preload("AdminBookmarkPoems").Where("id", id).First(&user); user.ID == 0 {
		return c.SendStatus(404)
	}
	var poem models.Poem
	if err := c.BodyParser(&poem); err != nil {
		return c.SendStatus(404)
	}

	err1 := database.DB.Model(&user).Association("AdminLikedPoems").Delete(&poem)
	if err1 != nil {
		return c.SendStatus(404)
	}

	database.DB.Save(&user)

	return c.JSON(user)
}

// GetLikedPoemsIdByAdminId returns only IDs of poems that user liked (optimized)
func GetLikedPoemsIdByAdminId(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)
	if userID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}

	var poemIDs []uint
	database.DB.Table("admin_liked_poems").
		Where("admin_id = ?", userID).
		Pluck("poem_id", &poemIDs)

	return c.JSON(poemIDs)
}

// GetLikedPoemsPaginated returns paginated list of poems that user liked
func GetLikedPoemsPaginated(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	userID := GetUserId(c)

	if uint(id) != userID {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// Get pagination parameters and search query
	params := helpers.GetPaginationParams(c)
	search := c.Query("search", "")

	// Get all liked poem IDs for this user
	var allPoemIDs []uint
	database.DB.Table("admin_liked_poems").
		Where("admin_id = ?", userID).
		Pluck("poem_id", &allPoemIDs)

	// If no liked poems, return empty result
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

	// Fetch paginated poems with like count
	var poems []PoemWithLikes
	database.DB.Table("poems").
		Select("poems.*, COUNT(admin_liked_poems.poem_id) as like_count").
		Joins("LEFT JOIN admin_liked_poems ON poems.id = admin_liked_poems.poem_id").
		Where("poems.id IN ? AND poems.is_deleted = ?", allPoemIDs, false).
		Group("poems.id").
		Offset(params.Offset).
		Limit(params.Limit).
		Scan(&poems)

	// Create paginated response
	response := helpers.CreatePaginationResponse(poems, total, params.Offset, params.Limit)

	fmt.Printf("[GetLikedPoemsPaginated] userID=%d, search=%s, total=%d, returned=%d\n", userID, search, total, len(poems))

	return c.JSON(response)
}
