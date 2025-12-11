package controllers

import (
	"backend/database"
	"backend/helpers"
	"backend/models"
	"backend/security"
	"backend/util"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

func CreateAdmin(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Validate and sanitize inputs
	validator := security.NewValidator()
	sanitizer := security.NewSanitizer()

	// Extract and validate username
	username, ok := data["username"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username is required",
		})
	}
	if err := validator.ValidateUsername(username); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	username = sanitizer.SanitizeString(username, 50)

	// Extract and validate email
	email, ok := data["email"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email is required",
		})
	}
	if err := validator.ValidateEmail(email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	email = sanitizer.SanitizeString(email, 255)

	// Extract and validate password
	password, ok := data["password"].(string)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Password is required",
		})
	}
	if err := validator.ValidatePassword(password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Parse and validate role_id
	var role uint
	switch v := data["role_id"].(type) {
	case string:
		roleInt, err := strconv.Atoi(v)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid role_id format",
			})
		}
		role = uint(roleInt)
	case float64:
		role = uint(v)
	case int:
		role = uint(v)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "role_id is required",
		})
	}

	if err := validator.ValidateRoleID(int(role)); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check if username already exists
	var existingAdmin models.Admin
	database.DB.Where("username = ?", username).First(&existingAdmin)
	if existingAdmin.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username already exists",
		})
	}

	// Check if email already exists
	database.DB.Where("email = ?", email).First(&existingAdmin)
	if existingAdmin.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Email already exists",
		})
	}

	// Hash password
	psw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}

	admin := models.Admin{
		Username: username,
		Email:    email,
		RoleID:   role,
		Password: psw,
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating admin",
		})
	}

	return c.JSON(GetAdminsBasicInfo())
}
func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate and sanitize inputs
	validator := security.NewValidator()
	sanitizer := security.NewSanitizer()

	// Validate username
	username := data["username"]
	if err := validator.ValidateUsername(username); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	username = sanitizer.SanitizeString(username, 50)

	// Validate email
	email := data["email"]
	if err := validator.ValidateEmail(email); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	email = sanitizer.SanitizeString(email, 255)

	// Validate password
	password := data["password"]
	if err := validator.ValidatePassword(password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check if username already exists
	var existingAdmin models.Admin
	database.DB.Where("username = ?", username).First(&existingAdmin)
	if existingAdmin.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bu kullanıcı adı daha önce alınmış",
		})
	}

	// Check if email already exists
	existingAdmin = models.Admin{}
	database.DB.Where("email = ?", email).First(&existingAdmin)
	if existingAdmin.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Bu e-posta adresi daha önce kullanılmış",
		})
	}

	// Hash password
	psw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error hashing password",
		})
	}

	// Create new user with role_id = 3
	admin := models.Admin{
		Username: username,
		Email:    email,
		RoleID:   3, // Default role for registered users
		Password: psw,
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"user": fiber.Map{
			"id":       admin.ID,
			"username": admin.Username,
			"email":    admin.Email,
			"role_id":  admin.RoleID,
		},
	})
}
func Login(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	// Validate and sanitize inputs
	validator := security.NewValidator()
	sanitizer := security.NewSanitizer()

	// Validate username
	username := data["username"]
	if err := validator.ValidateString("username", username, 1, 50, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	username = sanitizer.SanitizeString(username, 50)

	// Validate password
	password := data["password"]
	if err := validator.ValidateString("password", password, 1, 128, true); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Find user by username
	var admin models.Admin
	database.DB.Where("username = ?", username).First(&admin)

	if admin.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword(admin.Password, []byte(password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// Generate token
	token, err := util.SetToken(strconv.Itoa(int(admin.ID)))
	if err != nil {
		fmt.Println("Error creating token:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating authentication token",
		})
	}

	// Set secure cookie
	secureEnv := os.Getenv("COOKIE_SECURE")
	cookieSecure, _ := strconv.ParseBool(secureEnv)
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 5),
		HTTPOnly: true,
		Secure:   cookieSecure,
		SameSite: "None",
	}
	c.Cookie(&cookie)

	// Get user ID from token
	id, err := util.GetUserWithToken(cookie.Value)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	// Load user with relationships
	var admin1 models.Admin
	database.DB.
		Preload("AdminLikedPoems").
		Preload("AdminBookmarkPoems").
		Preload("UserBooksRead").
		Where("id = ?", id).
		Find(&admin1)

	if admin1.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Create log
	if err := CreateLog(c, admin1.Username, admin1.RoleID, admin1.ID); err != nil {
		fmt.Println("Log error:", err)
	}

	return c.JSON(fiber.Map{
		"admin":  admin1,
		"cookie": cookie,
	})
}

type StringData struct {
	Data string `json:"data"`
}

func User(c *fiber.Ctx) error {

	var data StringData
	if err := c.BodyParser(&data); err != nil {
		return c.SendString("Login veri yok.")
	}
	fmt.Println(data)
	cookie := data.Data
	id, err := util.GetUserWithToken(cookie)
	if err != nil {
		return c.SendStatus(401)
	}

	var admin models.Admin
	database.DB.Preload("AdminLikedPoems").Preload("AdminBookmarkPoems").Where("id", id).Find(&admin)

	if admin.ID == 0 {
		return c.SendStatus(405)
	}

	return c.JSON(fiber.Map{
		"admin": admin,
	})
}
func LogOut(c *fiber.Ctx) error {
	secureEnv := os.Getenv("COOKIE_SECURE")
	cookieSecure, _ := strconv.ParseBool(secureEnv)
	cookie := fiber.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-time.Hour), // geçmiş tarih, tarayıcı siler
		HTTPOnly: true,
		Secure:   cookieSecure, // prod'da true
		SameSite: "None",       // login ile aynı policy olmalı
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "Succsess",
	})
}

func DeleteAdmin(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var admin models.Admin
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	database.DB.Table("admins").Where("id", id).Delete(&admin)
	return c.JSON(GetAdminsBasicInfo())
}
func UpdateAdmin(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	id, _ := strconv.Atoi(c.Params("id"))
	var admin models.Admin
	if database.DB.Table("admins").Where("id", id).First(&admin); admin.ID == 0 {
		return c.SendStatus(404)
	}
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Parse role_id (can be string or number)
	var role uint
	switch v := data["role_id"].(type) {
	case string:
		roleInt, _ := strconv.Atoi(v)
		role = uint(roleInt)
	case float64:
		role = uint(v)
	case int:
		role = uint(v)
	}

	// Update fields
	admin.Username = data["username"].(string)
	admin.Email = data["email"].(string)
	admin.RoleID = role

	// Only update password if provided
	if password, ok := data["password"].(string); ok && password != "" {
		psw, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
		admin.Password = psw
	}

	database.DB.Model(&admin).Updates(admin)
	return c.JSON(GetAdminsBasicInfo())

}
func GetAdmin(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	var admin models.Admin
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	err1 := database.DB.Table("admins").Preload("AdminLikedPoems").Preload("AdminBookmarkPoems").Where("id", id).First(&admin).Error
	if err1 != nil {
		return c.SendStatus(404)
	}
	if admin.ID == 0 {
		return c.SendString("admin not found")
	}
	return c.JSON(admin)
}
func GetAdmins(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}
	return c.JSON(GetAllAdmins(c))
}

// AdminBasicInfo contains only necessary fields for admin management
type AdminBasicInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
}

// GetAdminsBasicInfo returns only basic info without relations
func GetAdminsBasicInfo() []AdminBasicInfo {
	var admins []AdminBasicInfo
	database.DB.Model(&models.Admin{}).Select("id, username, email, role_id").Find(&admins)
	return admins
}

// GetAdminsForManagement returns only necessary fields for admin management page
func GetAdminsForManagement(c *fiber.Ctx) error {
	roleID, err := helpers.GetUserRole(c)
	if err != nil || roleID != 1 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Access denied",
		})
	}

	return c.JSON(GetAdminsBasicInfo())
}
func GetAllAdmins(c *fiber.Ctx) []models.Admin {
	var admins []models.Admin
	result := database.DB.Preload("AdminLikedPoems").Preload("AdminBookmarkPoems").Find(&admins)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.Status(501)
		return nil
	}
	return admins
}
func AuthCheck(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "no token",
		})
	}

	userID, err := util.GetUserWithToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	var admin models.Admin
	result := database.DB.Select("id, username, role_id, profile_image, is_private").Where("id", userID).First(&admin)
	if result.Error != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "ok",
		"user": fiber.Map{
			"id":            admin.ID,
			"username":      admin.Username,
			"role_id":       admin.RoleID,
			"profile_image": admin.ProfileImage,
			"is_private":    admin.IsPrivate,
		},
	})
}
func GetUserId(c *fiber.Ctx) uint {
	// Get token from cookie
	cookie := c.Cookies("token")
	if cookie == "" {
		return uint(0)
	}

	// Extract user ID from token
	userID, err := util.GetUserWithToken(cookie)
	if err != nil {
		return uint(0)
	}

	// Convert userID string to uint
	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return uint(0)
	}
	return uint(id)

}

// UploadProfileImage handles profile image upload
func UploadProfileImage(c *fiber.Ctx) error {
	// Get user ID from token
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Get the file from the request
	file, err := c.FormFile("profile_image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Dosya yüklenemedi",
			"error":   err.Error(),
		})
	}

	// Validate file using FileValidator
	fileValidator := security.NewFileValidator()
	sanitizer := security.NewSanitizer()

	// Validate file size and magic bytes (max 1MB)
	maxSize := int64(1024 * 1024)
	if err := fileValidator.ValidateImageFile(file, maxSize); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Validate filename
	if err := fileValidator.ValidateFilename(file.Filename); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get safe extension
	ext, err := fileValidator.GetSafeExtension(file.Filename)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Sanitize filename
	sanitizedFilename := sanitizer.SanitizeFilename(file.Filename)

	// Generate unique filename
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("profile_%d_%d%s", userID, timestamp, ext)

	// Save the file
	uploadDir := "./uploads/profiles"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Klasör oluşturulamadı",
			"error":   err.Error(),
		})
	}

	filePath := fmt.Sprintf("%s/%s", uploadDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Dosya kaydedilemedi",
			"error":   err.Error(),
		})
	}

	// Update user's profile image in database
	var admin models.Admin
	if err := database.DB.Where("id = ?", userID).First(&admin).Error; err != nil {
		// Delete uploaded file if user not found
		os.Remove(filePath)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Kullanıcı bulunamadı",
		})
	}

	// Delete old profile image if exists
	if admin.ProfileImage != "" {
		oldPath := fmt.Sprintf(".%s", admin.ProfileImage)
		os.Remove(oldPath)
	}

	// Update profile image path
	profileImageURL := fmt.Sprintf("/uploads/profiles/%s", filename)
	admin.ProfileImage = profileImageURL

	if err := database.DB.Save(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Profil resmi güncellenemedi",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":       "Profil resmi başarıyla yüklendi",
		"profile_image": profileImageURL,
		"original_name": sanitizedFilename,
	})
}

// GetUserProfile returns user profile by username
func GetUserProfile(c *fiber.Ctx) error {
	username := c.Params("username")

	// Get current user ID (viewer) first
	viewerID := GetUserId(c)

	// First, load basic user info without relations
	var admin models.Admin
	if err := database.DB.
		Where("username = ?", username).
		First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Kullanıcı bulunamadı",
		})
	}

	// Check if viewer is friend with profile owner
	isFriend := AreFriends(viewerID, admin.ID)
	isOwnProfile := viewerID == admin.ID
	isPublic := !admin.IsPrivate

	// Check friendship status
	var friendshipStatus string = "none" // none, pending_sent, pending_received, accepted
	var friendship models.Friendship

	if !isOwnProfile {
		// Check if there's any friendship record
		err := database.DB.Where(
			"(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			viewerID, admin.ID, admin.ID, viewerID,
		).First(&friendship).Error

		if err == nil {
			// Friendship record exists
			if friendship.Status == "accepted" {
				friendshipStatus = "accepted"
			} else if friendship.Status == "pending" {
				// Check who sent the request
				if friendship.UserID == viewerID {
					friendshipStatus = "pending_sent"
				} else {
					friendshipStatus = "pending_received"
				}
			}
		}
	}

	// Security: Determine what data to load based on permissions
	canViewDetails := isOwnProfile || isFriend || isPublic

	// Get viewer's role_id for community filtering
	var viewer models.Admin
	viewerRoleID := uint(0) // Default to 0 (not logged in)
	if viewerID > 0 {
		database.DB.Select("role_id").Where("id = ?", viewerID).First(&viewer)
		viewerRoleID = viewer.RoleID
	}

	// Always return only counts for performance optimization
	// Items will be loaded lazily when user clicks on stats
	var likedCount, bookmarkCount, readBooksCount, commentsCount int64

	// Filter counts based on viewer's role_id and community
	// role_id 0 (not logged in) or 3 (Misafir): only show community=2 (public) content
	// role_id 1,2 (Admin/Kullanıcı): show all content
	if viewerRoleID == 0 || viewerRoleID == 3 {
		// Not logged in or Misafir: only count community=2 items
		database.DB.Table("admin_liked_poems").
			Joins("JOIN poems ON admin_liked_poems.poem_id = poems.id").
			Where("admin_liked_poems.admin_id = ? AND poems.community = ?", admin.ID, 2).
			Count(&likedCount)
		database.DB.Table("admin_bookmark_poems").
			Joins("JOIN poems ON admin_bookmark_poems.poem_id = poems.id").
			Where("admin_bookmark_poems.admin_id = ? AND poems.community = ?", admin.ID, 2).
			Count(&bookmarkCount)
		database.DB.Table("user_books_read").
			Joins("JOIN books ON user_books_read.book_id = books.id").
			Where("user_books_read.admin_id = ? AND books.community = ?", admin.ID, 2).
			Count(&readBooksCount)
	} else {
		// Admin/Kullanıcı: count all items
		database.DB.Table("admin_liked_poems").Where("admin_id = ?", admin.ID).Count(&likedCount)
		database.DB.Table("admin_bookmark_poems").Where("admin_id = ?", admin.ID).Count(&bookmarkCount)
		database.DB.Table("user_books_read").Where("admin_id = ?", admin.ID).Count(&readBooksCount)
	}

	database.DB.Model(&models.Comment{}).Where("admin_id = ?", admin.ID).Count(&commentsCount)

	// Security: Don't return password
	admin.Password = nil

	// Security: Don't return email unless it's own profile
	if !isOwnProfile {
		admin.Email = ""
	}

	response := fiber.Map{
		"user": fiber.Map{
			"id":            admin.ID,
			"username":      admin.Username,
			"role_id":       admin.RoleID,
			"is_private":    admin.IsPrivate,
			"profile_image": admin.ProfileImage,
			"email":         admin.Email,
			// Counts only - items will be loaded lazily
			"liked_poems_count":    likedCount,
			"bookmark_poems_count": bookmarkCount,
			"read_books_count":     readBooksCount,
			"comments_count":       commentsCount,
		},
		"is_friend":         isFriend,
		"is_own_profile":    isOwnProfile,
		"friendship_status": friendshipStatus,
		"can_view_details":  canViewDetails,
	}

	return c.JSON(response)
}

// UpdatePrivacy updates user's privacy setting
func UpdatePrivacy(c *fiber.Ctx) error {
	// Get userID from cookie - güvenlik için cookie'den alıyoruz
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	isPrivate, ok := data["is_private"].(bool)
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "is_private field is required and must be a boolean",
		})
	}

	// Sadece cookie'den gelen userID'nin kendi privacy ayarını güncelliyoruz
	// Bu sayede başka kullanıcının privacy ayarını değiştiremez
	if err := database.DB.Model(&models.Admin{}).
		Where("id = ?", userID).
		Update("is_private", isPrivate).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update privacy setting",
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Privacy setting updated successfully",
		"is_private": isPrivate,
	})
}

// GetUserLikedPoems returns user's liked poems (lazy loading)
func GetUserLikedPoems(c *fiber.Ctx) error {
	username := c.Params("username")

	// Get current user ID (viewer)
	viewerID := GetUserId(c)

	// Get target user
	var admin models.Admin
	if err := database.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Kullanıcı bulunamadı",
		})
	}

	// Check permissions
	isFriend := AreFriends(viewerID, admin.ID)
	isOwnProfile := viewerID == admin.ID
	isPublic := !admin.IsPrivate
	canViewDetails := isOwnProfile || isFriend || isPublic

	if !canViewDetails {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Bu içeriği görüntülemek için yetkiniz yok",
		})
	}

	// Get viewer's role_id for community filtering
	var viewer models.Admin
	viewerRoleID := uint(0) // Default to 0 (not logged in)
	if viewerID > 0 {
		database.DB.Select("role_id").Where("id = ?", viewerID).First(&viewer)
		viewerRoleID = viewer.RoleID
	}

	// Load liked poems with filtering based on role
	// First, get the poem IDs from the junction table
	var poemIDs []uint
	query := database.DB.
		Table("admin_liked_poems").
		Select("admin_liked_poems.poem_id").
		Joins("JOIN poems ON admin_liked_poems.poem_id = poems.id").
		Where("admin_liked_poems.admin_id = ?", admin.ID)

	// Filter by community based on viewer's role
	// role_id 0 (not logged in) or 3 (Misafir): only show community=2
	// role_id 1,2 (Admin/Kullanıcı): show all
	if viewerRoleID == 0 || viewerRoleID == 3 {
		query = query.Where("poems.community = ?", 2)
	}

	if err := query.Pluck("poem_id", &poemIDs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Şiirler yüklenirken hata oluştu",
		})
	}

	// Now load the full poems with author data
	var likedPoems []models.Poem
	if len(poemIDs) > 0 {
		if err := database.DB.
			Preload("AuthorData").
			Where("id IN ?", poemIDs).
			Find(&likedPoems).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Şiirler yüklenirken hata oluştu",
			})
		}
	}

	return c.JSON(fiber.Map{
		"poems": likedPoems,
	})
}

// GetUserReadBooks returns user's read books (lazy loading)
func GetUserReadBooks(c *fiber.Ctx) error {
	username := c.Params("username")

	// Get current user ID (viewer)
	viewerID := GetUserId(c)

	// Get target user
	var admin models.Admin
	if err := database.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Kullanıcı bulunamadı",
		})
	}

	// Check permissions
	isFriend := AreFriends(viewerID, admin.ID)
	isOwnProfile := viewerID == admin.ID
	isPublic := !admin.IsPrivate
	canViewDetails := isOwnProfile || isFriend || isPublic

	if !canViewDetails {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Bu içeriği görüntülemek için yetkiniz yok",
		})
	}

	// Get viewer's role_id for community filtering
	var viewer models.Admin
	viewerRoleID := uint(0) // Default to 0 (not logged in)
	if viewerID > 0 {
		database.DB.Select("role_id").Where("id = ?", viewerID).First(&viewer)
		viewerRoleID = viewer.RoleID
	}

	// Load read books with filtering based on role
	// First, get the book IDs from the junction table
	var bookIDs []uint
	query := database.DB.
		Table("user_books_read").
		Select("user_books_read.book_id").
		Joins("JOIN books ON user_books_read.book_id = books.id").
		Where("user_books_read.admin_id = ?", admin.ID)

	// Filter by community based on viewer's role
	// role_id 0 (not logged in) or 3 (Misafir): only show community=2
	// role_id 1,2 (Admin/Kullanıcı): show all
	if viewerRoleID == 0 || viewerRoleID == 3 {
		query = query.Where("books.community = ?", 2)
	}

	if err := query.Pluck("book_id", &bookIDs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Kitaplar yüklenirken hata oluştu",
		})
	}

	// Now load the full books with author data
	var readBooks []models.Book
	if len(bookIDs) > 0 {
		if err := database.DB.
			Preload("AuthorData").
			Where("id IN ? AND is_deleted = ?", bookIDs, false).
			Find(&readBooks).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Kitaplar yüklenirken hata oluştu",
			})
		}
	}

	return c.JSON(fiber.Map{
		"books": readBooks,
	})
}

// GetUserBookmarkedPoems retrieves all bookmarked poems for a specific user
func GetUserBookmarkedPoems(c *fiber.Ctx) error {
	username := c.Params("username")

	// Get current user ID (viewer)
	viewerID := GetUserId(c)

	// Get target user
	var admin models.Admin
	if err := database.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Kullanıcı bulunamadı",
		})
	}

	// Bookmarked poems are only visible to the owner
	if viewerID != admin.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Kaydedilen şiirleri sadece kendi profilinizde görüntüleyebilirsiniz",
		})
	}

	// Get viewer's role_id for community filtering
	var viewer models.Admin
	viewerRoleID := uint(0) // Default to 0 (not logged in)
	if viewerID > 0 {
		database.DB.Select("role_id").Where("id = ?", viewerID).First(&viewer)
		viewerRoleID = viewer.RoleID
	}

	// Load bookmarked poems with filtering based on role
	// First, get the poem IDs from the junction table
	var poemIDs []uint
	query := database.DB.
		Table("admin_bookmark_poems").
		Select("admin_bookmark_poems.poem_id").
		Joins("JOIN poems ON admin_bookmark_poems.poem_id = poems.id").
		Where("admin_bookmark_poems.admin_id = ?", admin.ID)

	// Filter by community based on viewer's role
	// role_id 0 (not logged in) or 3 (Misafir): only show community=2
	// role_id 1,2 (Admin/Kullanıcı): show all
	if viewerRoleID == 0 || viewerRoleID == 3 {
		query = query.Where("poems.community = ?", 2)
	}

	if err := query.Pluck("poem_id", &poemIDs).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Şiirler yüklenirken hata oluştu",
		})
	}

	// Now load the full poems with author data
	var bookmarkedPoems []models.Poem
	if len(poemIDs) > 0 {
		if err := database.DB.
			Preload("AuthorData").
			Where("id IN ?", poemIDs).
			Find(&bookmarkedPoems).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Şiirler yüklenirken hata oluştu",
			})
		}
	}

	return c.JSON(fiber.Map{
		"poems": bookmarkedPoems,
	})
}

// GetUserComments retrieves all comments made by a specific user
func GetUserComments(c *fiber.Ctx) error {
	username := c.Params("username")

	// Get current user ID (viewer)
	viewerID := GetUserId(c)

	// Get target user
	var admin models.Admin
	if err := database.DB.Where("username = ?", username).First(&admin).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Kullanıcı bulunamadı",
		})
	}

	// Comments are only visible to the owner
	if viewerID != admin.ID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Yorumlarınızı sadece kendi profilinizde görüntüleyebilirsiniz",
		})
	}

	// Get all comments by this user
	var comments []models.Comment
	if err := database.DB.
		Where("admin_id = ? AND is_deleted = ?", admin.ID, false).
		Preload("Book").
		Preload("Admin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email", "role_id", "profile_image")
		}).
		Order("id DESC").
		Find(&comments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Yorumlar yüklenirken hata oluştu",
		})
	}

	return c.JSON(fiber.Map{
		"comments": comments,
	})
}
