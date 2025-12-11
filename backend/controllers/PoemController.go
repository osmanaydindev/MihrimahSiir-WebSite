package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"backend/security"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"math"
	"strconv"
	"strings"
	"time"
)

// Helper function to get user's role_id from token
// func getUserRoleFromToken(c *fiber.Ctx) (uint, error) {
// 	cookie := c.Cookies("token")
// 	if cookie == "" {
// 		return 0, fiber.NewError(fiber.StatusUnauthorized, "No token provided")
// 	}

// 	userID, err := util.GetUserWithToken(cookie)
// 	if err != nil {
// 		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
// 	}

// 	var admin models.Admin
// 	result := database.DB.Where("id = ?", userID).First(&admin)
// 	if result.Error != nil {
// 		return 0, fiber.NewError(fiber.StatusUnauthorized, "User not found")
// 	}

// 	return admin.RoleID, nil
// }

// Helper function to apply community filter based on role_id
func applyCommunityFilter(db *gorm.DB, roleID uint) *gorm.DB {
	// role_id 1 and 2 can see all poems (community 1 and 2)
	// role_id 3 can only see public poems (community 2)
	if roleID == 1 || roleID == 2 {
		return db // No filter, can see all
	}
	return db.Where("community = ?", 2) // Only public poems
}

func replaceChars(s string, replacements map[rune]rune) string {
	var result strings.Builder
	for _, char := range s {
		if replacement, ok := replacements[char]; ok {
			result.WriteRune(replacement)
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}

func CreatePoem(c *fiber.Ctx) error {
	var poem models.Poem
	if err := c.BodyParser(&poem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate and sanitize inputs
	validator := security.NewValidator()
	sanitizer := security.NewSanitizer()

	// Validate title
	if err := validator.ValidateString("title", poem.Title, 1, 255, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate content
	if err := validator.ValidateString("content", poem.Content, 1, 50000, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate community
	if err := validator.ValidateCommunity(poem.Community); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Validate author_id if provided
	if poem.AuthorID != nil {
		if err := validator.ValidateID("author_id", *poem.AuthorID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Sanitize inputs to prevent XSS
	poem.Title = sanitizer.SanitizeString(poem.Title, 255)
	poem.Content = sanitizer.SanitizeString(poem.Content, 50000)

	// Check for dangerous content
	if sanitizer.ContainsDangerousContent(poem.Title) || sanitizer.ContainsDangerousContent(poem.Content) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content contains potentially dangerous elements",
		})
	}

	poem.IsDeleted = false
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
	slug := replaceChars(poem.Title, x)
	poem.CreatedAt = time.Now().Format("02-01-2006")
	poem.CreatedAtParse = time.Now().String()
	poem.Slug = strings.ToLower(slug)

	if err := database.DB.Create(&poem).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create poem",
		})
	}

	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getPoems(roleID))
}
func DeletePoem(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var poem models.Poem
	database.DB.Table("poems").Where("id", id).Find(&poem)
	poem.IsDeleted = true
	database.DB.Save(&poem)

	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getPoems(roleID))
}
func UpdatePoem(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid poem ID",
		})
	}

	var poem models.Poem
	result := database.DB.Table("poems").Where("id", id).Find(&poem)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Poem not found",
		})
	}

	var updateData models.Poem
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate and sanitize inputs
	validator := security.NewValidator()
	sanitizer := security.NewSanitizer()

	// Validate title if provided
	if updateData.Title != "" {
		if err := validator.ValidateString("title", updateData.Title, 1, 255, true); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		updateData.Title = sanitizer.SanitizeString(updateData.Title, 255)
	}

	// Validate content if provided
	if updateData.Content != "" {
		if err := validator.ValidateString("content", updateData.Content, 1, 50000, true); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		updateData.Content = sanitizer.SanitizeString(updateData.Content, 50000)
	}

	// Validate community if provided
	if updateData.Community != 0 {
		if err := validator.ValidateCommunity(updateData.Community); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
	}

	// Check for dangerous content
	if sanitizer.ContainsDangerousContent(updateData.Title) || sanitizer.ContainsDangerousContent(updateData.Content) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content contains potentially dangerous elements",
		})
	}

	if err := database.DB.Model(&poem).Updates(updateData).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update poem",
		})
	}

	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getPoems(roleID))
}
func GetPoem(c *fiber.Ctx) error {
	slug := c.Params("slug")
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	poem := models.Poem{
		Slug: slug,
	}

	// Apply community filter when fetching the main poem
	query := database.DB.Table("poems").Where("slug", slug)
	query = applyCommunityFilter(query, roleID)
	result := query.Preload("AuthorData").First(&poem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.SendString("poem not found")
	}

	// Apply community filter to random poems as well
	var randomPoem []models.Poem
	randomQuery := database.DB.
		Where("id != ?", poem.ID).
		Where("is_deleted != ?", true)
	randomQuery = applyCommunityFilter(randomQuery, roleID)
	if err := randomQuery.Preload("AuthorData").
		Order("RANDOM()").
		Limit(2).
		Find(&randomPoem).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"poem":       poem,
		"randomPoem": randomPoem,
	})
}
func GetPoemById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	poem := models.Poem{
		ID: uint(id),
	}

	query := database.DB.Table("poems").Where("id", id)
	query = applyCommunityFilter(query, roleID)
	result := query.Preload("AuthorData").First(&poem)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.SendString("poem not found")
	}
	return c.JSON(poem)
}
func GetAllPoems(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	return c.JSON(getPoems(roleID))
}
func getPoems(roleID uint) []models.Poem {
	var poems []models.Poem
	query := database.DB.Where("is_deleted", false)
	query = applyCommunityFilter(query, roleID)
	query.Preload("AuthorData").Order("created_at desc").Find(&poems)
	return poems
}
func divideAndRoundUp(a, b int) int {
	if b == 0 {
		fmt.Println("Division by zero is not allowed")
		return 0
	}

	result := float64(a) / float64(b)
	roundedResult := math.Ceil(result) // Sonucu yukarı yuvarlar
	return int(roundedResult)
}
func GetLatestPoems(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	var poems []PoemWithLikes

	query := database.DB.Table("poems").
		Select("poems.*, COUNT(admin_liked_poems.poem_id) as like_count").
		Joins("LEFT JOIN admin_liked_poems ON poems.id = admin_liked_poems.poem_id").
		Where("poems.is_deleted = ?", false).
		Group("poems.id").
		Order("poems.created_at_parse desc").
		Limit(5)

	query = applyCommunityFilter(query, roleID)
	query.Scan(&poems)

	return c.JSON(poems)
}
func GetPoemsPaginated(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 10
	offset := (page - 1) * limit
	var total int64
	var poems []PoemWithLikes

	query := database.DB.Table("poems").
		Select("poems.*, COUNT(admin_liked_poems.poem_id) as like_count").
		Joins("LEFT JOIN admin_liked_poems ON poems.id = admin_liked_poems.poem_id").
		Where("poems.is_deleted = ?", false).
		Group("poems.id").
		Offset(offset).
		Limit(limit)

	query = applyCommunityFilter(query, roleID)
	query.Scan(&poems)

	countQuery := database.DB.Model(&models.Poem{}).Where("is_deleted", false)
	countQuery = applyCommunityFilter(countQuery, roleID)
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"poems": poems,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": divideAndRoundUp(int(total), limit),
		},
	})
}
func GetSearchPoems(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	// Validate and sanitize inputs
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}
	if page > 10000 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Page number too large",
		})
	}

	search := c.Query("search")

	// Sanitize search input to prevent ReDoS and SQL injection
	sanitizer := security.NewSanitizer()
	validator := security.NewValidator()

	// Validate search length
	if err := validator.ValidateString("search", search, 0, 200, false); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Sanitize search for SQL LIKE query (prevents SQL injection)
	sanitizedSearch := sanitizer.SanitizeSQLLike(search)

	// Check for dangerous content
	if sanitizer.ContainsDangerousContent(search) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search contains invalid characters",
		})
	}

	limit := 10
	offset := (page - 1) * limit
	var total int64
	var poems []PoemWithLikes

	// Use ILIKE instead of regex to prevent ReDoS vulnerability
	// ILIKE is PostgreSQL's case-insensitive LIKE
	searchPattern := "%" + sanitizedSearch + "%"

	// Apply community filter in raw SQL query with like count
	var rawQuery string
	if roleID == 1 || roleID == 2 {
		// Admin/Member: can see all poems
		rawQuery = `SELECT poems.*, COUNT(admin_liked_poems.poem_id) as like_count
					FROM poems
					LEFT JOIN admin_liked_poems ON poems.id = admin_liked_poems.poem_id
					WHERE (title ILIKE ? OR content ILIKE ?) AND is_deleted = false
					GROUP BY poems.id
					ORDER BY poems.created_at_parse DESC
					OFFSET ? LIMIT ?`
		database.DB.Raw(rawQuery, searchPattern, searchPattern, offset, limit).Scan(&poems)
	} else {
		// Guest (role_id 3): can only see public poems (community = 2)
		rawQuery = `SELECT poems.*, COUNT(admin_liked_poems.poem_id) as like_count
					FROM poems
					LEFT JOIN admin_liked_poems ON poems.id = admin_liked_poems.poem_id
					WHERE (title ILIKE ? OR content ILIKE ?) AND is_deleted = false AND community = 2
					GROUP BY poems.id
					ORDER BY poems.created_at_parse DESC
					OFFSET ? LIMIT ?`
		database.DB.Raw(rawQuery, searchPattern, searchPattern, offset, limit).Scan(&poems)
	}

	// Count query with community filter
	countQuery := database.DB.Model(&models.Poem{}).
		Where("title ILIKE ? OR content ILIKE ?", searchPattern, searchPattern).
		Where("is_deleted", false)
	countQuery = applyCommunityFilter(countQuery, roleID)
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"poems": poems,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": divideAndRoundUp(int(total), limit),
		},
	})
}

// PoemWithLikes - Poem struct with like count (without gorm:"-" tag)
type PoemWithLikes struct {
	ID             uint   `json:"id"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	AuthorID       *uint  `json:"author_id"`
	Content        string `json:"content"`
	IsDeleted      bool   `json:"is_deleted"`
	Slug           string `json:"slug"`
	CreatedAt      string `json:"created_at"`
	CreatedAtParse string `json:"created_at_parse"`
	Community      int    `json:"community"`
	LikeCount      int    `json:"like_count"`
}

// GetPopularPoems - Get poems ordered by like count
func GetPopularPoems(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil {
		return err
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit := 10
	offset := (page - 1) * limit
	var total int64
	var poems []PoemWithLikes

	// Query poems with like count
	query := database.DB.Table("poems").
		Select("poems.*, COUNT(admin_liked_poems.poem_id) as like_count").
		Joins("LEFT JOIN admin_liked_poems ON poems.id = admin_liked_poems.poem_id").
		Where("poems.is_deleted = ?", false).
		Group("poems.id").
		Order("like_count DESC")

	query = applyCommunityFilter(query, roleID)
	query.Offset(offset).Limit(limit).Scan(&poems)

	// Count query
	countQuery := database.DB.Model(&models.Poem{}).Where("is_deleted", false)
	countQuery = applyCommunityFilter(countQuery, roleID)
	countQuery.Count(&total)

	return c.JSON(fiber.Map{
		"poems": poems,
		"meta": fiber.Map{
			"total":     total,
			"page":      page,
			"last_page": divideAndRoundUp(int(total), limit),
		},
	})
}
