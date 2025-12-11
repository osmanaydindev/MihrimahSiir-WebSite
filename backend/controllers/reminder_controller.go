package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

// CreateReminder yeni bir hatırlatıcı oluşturur
func CreateReminder(c *fiber.Ctx) error {
	var newReminder models.Reminder
	if err := c.BodyParser(&newReminder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Hatırlatıcıyı veritabanına kaydet
	if err := database.DB.Create(&newReminder).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create reminder"})
	}

	return c.JSON(GetAllReminders(c))
}

// GetReminder belirli bir hatırlatıcıyı getirir
func GetReminder(c *fiber.Ctx) error {
	id := c.Params("id")
	var reminder models.Reminder

	// Veritabanından hatırlatıcıyı bul
	if err := database.DB.First(&reminder, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reminder not found"})
	}
	return c.JSON(reminder)
}

func GetAllReminders(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	var reminders []models.Reminder

	if roleID == 1 {
		database.DB.Where("permission IN ?", []int{1, 2, 3}).Find(&reminders)
	} else {
		database.DB.Where("permission = ?", roleID).Find(&reminders)
	}

	return c.JSON(reminders)
}

// GetRemindersPaginated returns paginated list of reminders with search support
func GetRemindersPaginated(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	// Get pagination parameters and search query
	params := helpers.GetPaginationParams(c)
	search := c.Query("search", "")

	// Build base query with permission filter
	var baseQuery = database.DB.Model(&models.Reminder{})
	if roleID == 1 {
		baseQuery = baseQuery.Where("permission IN ?", []int{1, 2, 3})
	} else {
		baseQuery = baseQuery.Where("permission = ?", roleID)
	}

	// Apply search filter if provided
	if search != "" {
		searchPattern := "%" + search + "%"
		baseQuery = baseQuery.Where("text ILIKE ?", searchPattern)
	}

	// Count total matching reminders
	var total int64
	baseQuery.Count(&total)

	// Fetch paginated reminders
	reminders := []models.Reminder{}
	query := database.DB.Model(&models.Reminder{})
	if roleID == 1 {
		query = query.Where("permission IN ?", []int{1, 2, 3})
	} else {
		query = query.Where("permission = ?", roleID)
	}

	// Apply search filter
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("text ILIKE ?", searchPattern)
	}

	query.Order("created_at DESC").
		Offset(params.Offset).
		Limit(params.Limit).
		Find(&reminders)

	// Create paginated response
	response := helpers.CreatePaginationResponse(reminders, total, params.Offset, params.Limit)

	fmt.Printf("[GetRemindersPaginated] roleID=%d, search=%s, total=%d, returned=%d\n", roleID, search, total, len(reminders))

	return c.JSON(response)
}

// UpdateReminder belirli bir hatırlatıcıyı günceller
func UpdateReminder(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedReminder models.Reminder
	if err := c.BodyParser(&updatedReminder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Veritabanında hatırlatıcıyı güncelle
	if err := database.DB.Model(&models.Reminder{}).Where("id = ?", id).Updates(updatedReminder).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reminder not found"})
	}
	return c.JSON(GetAllReminders(c))
}

// DeleteReminder belirli bir hatırlatıcıyı siler
func DeleteReminder(c *fiber.Ctx) error {
	id := c.Params("id")
	var reminder models.Reminder

	// Veritabanından hatırlatıcıyı sil
	if err := database.DB.Delete(&reminder, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Reminder not found"})
	}

	return c.JSON(GetAllReminders(c))
}
