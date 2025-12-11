package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func AddComment(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	bookId, _ := strconv.Atoi(data["book_id"])
	adminId, _ := strconv.Atoi(data["admin_id"])

	// Page is optional, can be nil
	var page *int
	if pageStr, ok := data["page"]; ok && pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil {
			page = &pageNum
		}
	}

	userID := GetUserId(c)
	if userID != uint(adminId) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	comment := models.Comment{
		AdminID:   uint(adminId),
		BookID:    uint(bookId),
		Title:     data["title"],
		Content:   data["content"],
		Page:      page,
		IsDeleted: false,
	}
	if err := database.DB.Create(&comment).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Hatalı İstek",
			"error":   err.Error(),
		})
	}
	comments := getCommentsOfBookForUser(int64(comment.BookID), userID)
	return c.JSON(comments)
}
func GetComments(c *fiber.Ctx) error {
	bookId, _ := strconv.Atoi(c.Params("book_id"))
	userID := GetUserId(c)

	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var comments []models.Comment

	// Admin can see all comments
	if roleID == 1 {
		database.DB.Where("book_id = ?", bookId).Preload("Admin").Preload("Book").Find(&comments)
		return c.JSON(comments)
	}

	// Regular users can only see comments from their friends (and themselves)
	friendIDs := GetFriendIDs(userID)

	database.DB.Where("book_id = ? AND admin_id IN ?", bookId, friendIDs).
		Preload("Admin").
		Preload("Book").
		Find(&comments)

	return c.JSON(comments)
}
func DeleteComment(c *fiber.Ctx) error {
	commentId, _ := strconv.Atoi(c.Params("comment_id"))
	var comment models.Comment
	database.DB.First(&comment, commentId)
	userID := GetUserId(c)
	if userID != comment.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	bookId := comment.BookID
	database.DB.Delete(&comment)
	comments := getCommentsOfBookForUser(int64(bookId), userID)
	return c.JSON(comments)
}
func UpdateComment(c *fiber.Ctx) error {
	commentId, _ := strconv.Atoi(c.Params("comment_id"))
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var comment models.Comment
	database.DB.First(&comment, commentId)

	// Check ownership - only the comment owner or admin can update
	if roleID != 1 && userID != comment.AdminID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}

	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	comment.Title = data["title"]
	comment.Content = data["content"]
	database.DB.Save(&comment)
	comments := getCommentsOfBookForUser(int64(comment.BookID), userID)
	return c.JSON(comments)
}

// getCommentsOfBook returns all comments for a book (admin only)
func getCommentsOfBook(id int64) []models.Comment {
	var comments []models.Comment
	database.DB.Where("book_id = ?", id).Preload("Admin").Preload("Book").Find(&comments)
	return comments
}

// getCommentsOfBookForUser returns comments filtered by friendship
func getCommentsOfBookForUser(bookID int64, userID uint) []models.Comment {
	var comments []models.Comment

	// Get user role
	var admin models.Admin
	database.DB.First(&admin, userID)

	// Admin can see all comments
	if admin.RoleID == 1 {
		database.DB.Where("book_id = ?", bookID).Preload("Admin").Preload("Book").Find(&comments)
		return comments
	}

	// Regular users can only see comments from their friends (and themselves)
	friendIDs := GetFriendIDs(userID)

	database.DB.Where("book_id = ? AND admin_id IN ?", bookID, friendIDs).
		Preload("Admin").
		Preload("Book").
		Find(&comments)

	return comments
}
