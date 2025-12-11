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

// CreateAuthor creates a new author
func CreateAuthor(c *fiber.Ctx) error {
	var author models.Author
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Generate slug from name
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
	slug := replaceChars(author.Name, x)
	author.Slug = strings.ToLower(slug)
	author.CreatedAt = time.Now().Format("02-01-2006")
	author.IsDeleted = false

	if err := database.DB.Create(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create author",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(author)
}

// GetAuthors returns paginated list of authors
func GetAuthors(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	// Get pagination parameters
	params := helpers.GetPaginationParams(c)

	// Get total count
	var total int64
	database.DB.Model(&models.Author{}).Where("is_deleted = ?", false).Count(&total)

	// Get paginated authors with their poems and books
	var authors []models.Author
	query := database.DB.Where("is_deleted = ?", false)

	// Preload poems and books with community filtering
	query = query.Preload("Poems", func(db *gorm.DB) *gorm.DB {
		q := db.Where("is_deleted = ?", false)
		if roleID != 1 && roleID != 2 {
			q = q.Where("community = ?", 2)
		}
		return q
	}).Preload("Books", func(db *gorm.DB) *gorm.DB {
		q := db.Where("is_deleted = ?", false)
		if roleID != 1 && roleID != 2 {
			q = q.Where("community = ?", 2)
		}
		return q
	})

	query.Offset(params.Offset).
		Limit(params.Limit).
		Order("name ASC").
		Find(&authors)

	// Create paginated response
	response := helpers.CreatePaginationResponse(authors, total, params.Offset, params.Limit)

	// Log for debugging
	fmt.Printf("[GetAuthors] roleID=%d, total=%d, returned=%d\n", roleID, total, len(authors))

	return c.JSON(response)
}

// GetAuthor returns a single author by slug with their poems and books
func GetAuthor(c *fiber.Ctx) error {
	slug := c.Params("slug")
	userID := GetUserId(c)
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	var author models.Author
	query := database.DB.Where("slug = ? AND is_deleted = ?", slug, false)

	// Preload poems and books with community filtering
	query = query.Preload("Poems", func(db *gorm.DB) *gorm.DB {
		q := db.Where("is_deleted = ?", false)
		if roleID != 1 && roleID != 2 {
			q = q.Where("community = ?", 2)
		}
		return q.Order("created_at DESC")
	}).Preload("Books", func(db *gorm.DB) *gorm.DB {
		q := db.Where("is_deleted = ?", false)
		if roleID != 1 && roleID != 2 {
			q = q.Where("community = ?", 2)
		}
		return q.Order("created_at DESC")
	})

	if err := query.First(&author).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	fmt.Printf("[GetAuthor] slug=%s, userID=%d, roleID=%d, poems=%d, books=%d\n",
		slug, userID, roleID, len(author.Poems), len(author.Books))

	return c.JSON(author)
}

// UpdateAuthor updates an existing author
func UpdateAuthor(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var author models.Author
	if err := database.DB.Where("id = ?", id).First(&author).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	var updateData models.Author
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	if updateData.Name != "" {
		author.Name = updateData.Name
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
		slug := replaceChars(author.Name, x)
		author.Slug = strings.ToLower(slug)
	}
	if updateData.Bio != "" {
		author.Bio = updateData.Bio
	}
	if updateData.BirthYear != nil {
		author.BirthYear = updateData.BirthYear
	}
	if updateData.DeathYear != nil {
		author.DeathYear = updateData.DeathYear
	}
	if updateData.Nationality != "" {
		author.Nationality = updateData.Nationality
	}
	if updateData.Image != "" {
		author.Image = updateData.Image
	}

	if err := database.DB.Save(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update author",
		})
	}

	return c.JSON(author)
}

// DeleteAuthor soft deletes an author
func DeleteAuthor(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var author models.Author
	if err := database.DB.Where("id = ?", id).First(&author).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	author.IsDeleted = true
	if err := database.DB.Save(&author).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete author",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Author deleted successfully",
	})
}

// GetAllAuthorsForDropdown returns all authors for dropdown selection (no pagination)
// Returns all fields for admin management
func GetAllAuthorsForDropdown(c *fiber.Ctx) error {
	var authors []models.Author
	database.DB.
		Where("is_deleted = ?", false).
		Order("name ASC").
		Find(&authors)

	return c.JSON(authors)
}

// GetAuthorById returns a single author by ID
func GetAuthorById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var author models.Author
	if err := database.DB.Where("id = ? AND is_deleted = ?", id, false).First(&author).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Author not found",
		})
	}

	return c.JSON(author)
}
