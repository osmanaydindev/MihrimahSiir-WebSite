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

// GetHomepageItems - Get homepage items based on user's role_id
func GetHomepageItems(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	var homepages []models.Homepage
	// Users can see homepage items ONLY for their exact role_id
	database.DB.Where("permission = ?", roleID).Find(&homepages)

	return c.JSON(homepages)
}

// GetAllHomepageItems - Get all homepage items (for admin panel)
func GetAllHomepageItems(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}

	var homepages []models.Homepage
	database.DB.Find(&homepages)

	return c.JSON(homepages)
}

// GetHomepageItem - Get single homepage item by ID
func GetHomepageItem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	var homepage models.Homepage
	result := database.DB.First(&homepage, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"message": "Homepage item not found",
		})
	}

	return c.JSON(homepage)
}

// CreateHomepageItem - Create new homepage item
func CreateHomepageItem(c *fiber.Ctx) error {
	var homepage models.Homepage

	if err := c.BodyParser(&homepage); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	database.DB.Create(&homepage)

	var homepages []models.Homepage
	database.DB.Find(&homepages)

	return c.JSON(homepages)
}

// UpdateHomepageItem - Update existing homepage item
func UpdateHomepageItem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	var homepage models.Homepage
	result := database.DB.First(&homepage, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"message": "Homepage item not found",
		})
	}

	if err := c.BodyParser(&homepage); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	database.DB.Model(&homepage).Updates(homepage)

	var homepages []models.Homepage
	database.DB.Find(&homepages)

	return c.JSON(homepages)
}

// DeleteHomepageItem - Delete homepage item
func DeleteHomepageItem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID",
		})
	}

	var homepage models.Homepage
	result := database.DB.First(&homepage, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(404).JSON(fiber.Map{
			"message": "Homepage item not found",
		})
	}

	database.DB.Delete(&homepage)

	var homepages []models.Homepage
	database.DB.Find(&homepages)

	return c.JSON(homepages)
}
