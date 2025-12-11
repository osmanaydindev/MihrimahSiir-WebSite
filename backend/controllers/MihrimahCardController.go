package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// GetAllMihrimahCards - Get all mihrimah cards
func GetAllMihrimahCards(c *fiber.Ctx) error {
	var cards []models.MihrimahCard
	database.DB.Find(&cards)
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID == 3 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	return c.JSON(cards)
}

// GetMihrimahCard - Get single mihrimah card by ID
func GetMihrimahCard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID != uint(id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	var card models.MihrimahCard
	result := database.DB.First(&card, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"message": "Mihrimah card not found",
		})
	}

	return c.JSON(card)
}

// CreateMihrimahCard - Create new mihrimah card
func CreateMihrimahCard(c *fiber.Ctx) error {
	var card models.MihrimahCard

	if err := c.BodyParser(&card); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	database.DB.Create(&card)

	var cards []models.MihrimahCard
	database.DB.Find(&cards)

	return c.JSON(cards)
}

// UpdateMihrimahCard - Update existing mihrimah card
func UpdateMihrimahCard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	var card models.MihrimahCard
	result := database.DB.First(&card, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"message": "Mihrimah card not found",
		})
	}

	if err := c.BodyParser(&card); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	database.DB.Model(&card).Updates(card)

	var cards []models.MihrimahCard
	database.DB.Find(&cards)

	return c.JSON(cards)
}

// DeleteMihrimahCard - Delete mihrimah card
func DeleteMihrimahCard(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	var card models.MihrimahCard
	result := database.DB.First(&card, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"message": "Mihrimah card not found",
		})
	}

	database.DB.Delete(&card)

	var cards []models.MihrimahCard
	database.DB.Find(&cards)

	return c.JSON(cards)
}
