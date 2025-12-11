package controllers

import (
	"backend/database"
	"backend/models"
	ws "backend/websocket"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// SendFriendRequest sends a friend request to another user by username
func SendFriendRequest(c *fiber.Ctx) error {
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

	friendUsername, ok := data["username"].(string)
	if !ok || friendUsername == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Username is required",
		})
	}

	// Find friend by username
	var friend models.Admin
	if err := database.DB.Where("username = ?", friendUsername).First(&friend).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// Check if trying to add themselves
	if friend.ID == userID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "You cannot add yourself as a friend",
		})
	}

	// Check if friendship already exists
	var existingFriendship models.Friendship
	err := database.DB.Where(
		"(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID, friend.ID, friend.ID, userID,
	).First(&existingFriendship).Error

	if err == nil {
		// Friendship exists
		if existingFriendship.Status == "accepted" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "You are already friends",
			})
		} else if existingFriendship.Status == "pending" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Friend request already sent",
			})
		}
	}

	// Create new friend request
	friendship := models.Friendship{
		UserID:   userID,
		FriendID: friend.ID,
		Status:   "pending",
	}

	if err := database.DB.Create(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to send friend request",
		})
	}

	// Send WebSocket notification to the friend
	count := getPendingRequestCount(friend.ID)
	fmt.Printf("[WebSocket] Sending friend_request_received to user %d, count: %d\n", friend.ID, count)
	ws.GlobalHub.SendToUser(friend.ID, "friend_request_received", map[string]interface{}{
		"count": count,
	})

	return c.JSON(fiber.Map{
		"message":    "Friend request sent successfully",
		"friendship": friendship,
	})
}

// GetFriendRequests retrieves pending friend requests received by the user
func GetFriendRequests(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var friendships []models.Friendship
	err := database.DB.
		Where("friend_id = ? AND status = ?", userID, "pending").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email", "role_id", "profile_image")
		}).
		Order("created_at DESC").
		Find(&friendships).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch friend requests",
		})
	}

	return c.JSON(friendships)
}

// GetSentRequests retrieves pending friend requests sent by the user
func GetSentRequests(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var friendships []models.Friendship
	err := database.DB.
		Where("user_id = ? AND status = ?", userID, "pending").
		Preload("Friend", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "username", "email", "role_id", "profile_image")
		}).
		Order("created_at DESC").
		Find(&friendships).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch sent requests",
		})
	}

	return c.JSON(friendships)
}

// GetFriends retrieves all accepted friends
func GetFriends(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var friendships []models.Friendship
	err := database.DB.
		Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, "accepted").
		Preload("User").
		Preload("Friend").
		Order("updated_at DESC").
		Find(&friendships).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch friends",
		})
	}

	// Format response to return friend info based on current user
	type FriendInfo struct {
		FriendshipID uint      `json:"friendship_id"`
		UserID       uint      `json:"user_id"`
		Username     string    `json:"username"`
		Email        string    `json:"email"`
		RoleID       uint      `json:"role_id"`
		ProfileImage string    `json:"profile_image"`
		SinceDate    time.Time `json:"since_date"`
	}

	var friends []FriendInfo
	for _, fs := range friendships {
		var friend FriendInfo
		friend.FriendshipID = fs.ID
		friend.SinceDate = fs.UpdatedAt

		if fs.UserID == userID {
			// Friend is the FriendID
			friend.UserID = fs.Friend.ID
			friend.Username = fs.Friend.Username
			friend.Email = fs.Friend.Email
			friend.RoleID = fs.Friend.RoleID
			friend.ProfileImage = fs.Friend.ProfileImage
		} else {
			// Friend is the UserID
			friend.UserID = fs.User.ID
			friend.Username = fs.User.Username
			friend.Email = fs.User.Email
			friend.RoleID = fs.User.RoleID
			friend.ProfileImage = fs.User.ProfileImage
		}
		friends = append(friends, friend)
	}

	return c.JSON(friends)
}

// AcceptFriendRequest accepts a pending friend request
func AcceptFriendRequest(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	requestID, _ := strconv.Atoi(c.Params("id"))

	var friendship models.Friendship
	err := database.DB.First(&friendship, requestID).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Friend request not found",
		})
	}

	// Check if the request is for this user
	if friendship.FriendID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You can only accept requests sent to you",
		})
	}

	// Check if already accepted
	if friendship.Status == "accepted" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Friend request already accepted",
		})
	}

	// Update status to accepted
	friendship.Status = "accepted"
	friendship.UpdatedAt = time.Now()

	if err := database.DB.Save(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to accept friend request",
		})
	}

	// Send WebSocket notification to the requester
	ws.GlobalHub.SendToUser(friendship.UserID, "friend_request_accepted", map[string]interface{}{
		"username": getUsernameByID(userID),
	})

	// Update notification count for current user
	ws.GlobalHub.SendToUser(userID, "friend_request_update", map[string]interface{}{
		"count": getPendingRequestCount(userID),
	})

	return c.JSON(fiber.Map{
		"message":    "Friend request accepted",
		"friendship": friendship,
	})
}

// RejectFriendRequest rejects a pending friend request
func RejectFriendRequest(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	requestID, _ := strconv.Atoi(c.Params("id"))

	var friendship models.Friendship
	err := database.DB.First(&friendship, requestID).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Friend request not found",
		})
	}

	// Check if the request is for this user
	if friendship.FriendID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You can only reject requests sent to you",
		})
	}

	// Delete the friendship request
	if err := database.DB.Delete(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to reject friend request",
		})
	}

	// Update notification count for current user (who rejected)
	ws.GlobalHub.SendToUser(userID, "friend_request_update", map[string]interface{}{
		"count": getPendingRequestCount(userID),
	})

	// Notify the requester that their request was rejected
	ws.GlobalHub.SendToUser(friendship.UserID, "friend_request_update", map[string]interface{}{
		"count": getPendingRequestCount(friendship.UserID),
	})

	return c.JSON(fiber.Map{
		"message": "Friend request rejected",
	})
}

// CancelFriendRequest cancels a sent friend request
func CancelFriendRequest(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	requestID, _ := strconv.Atoi(c.Params("id"))

	var friendship models.Friendship
	err := database.DB.First(&friendship, requestID).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Friend request not found",
		})
	}

	// Check if the request was sent by this user
	if friendship.UserID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You can only cancel requests you sent",
		})
	}

	// Check if still pending
	if friendship.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot cancel non-pending request",
		})
	}

	// Delete the friendship request
	if err := database.DB.Delete(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to cancel friend request",
		})
	}

	// Update notification count for the friend (who would have received it)
	ws.GlobalHub.SendToUser(friendship.FriendID, "friend_request_update", map[string]interface{}{
		"count": getPendingRequestCount(friendship.FriendID),
	})

	// Notify the user who cancelled that their sent request list should update
	ws.GlobalHub.SendToUser(userID, "friend_request_update", map[string]interface{}{
		"count": getPendingRequestCount(userID),
	})

	return c.JSON(fiber.Map{
		"message": "Friend request cancelled",
	})
}

// RemoveFriend removes an accepted friendship
func RemoveFriend(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	friendshipID, _ := strconv.Atoi(c.Params("id"))

	var friendship models.Friendship
	err := database.DB.First(&friendship, friendshipID).Error
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Friendship not found",
		})
	}

	// Check if user is part of this friendship
	if friendship.UserID != userID && friendship.FriendID != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "You can only remove your own friendships",
		})
	}

	// Determine the other user in the friendship
	var otherUserID uint
	if friendship.UserID == userID {
		otherUserID = friendship.FriendID
	} else {
		otherUserID = friendship.UserID
	}

	// Delete the friendship
	if err := database.DB.Delete(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to remove friend",
		})
	}

	// Notify both users that the friendship was removed
	ws.GlobalHub.SendToUser(userID, "friend_removed", map[string]interface{}{
		"message": "Friendship removed",
	})
	ws.GlobalHub.SendToUser(otherUserID, "friend_removed", map[string]interface{}{
		"message": "Friendship removed",
	})

	return c.JSON(fiber.Map{
		"message": "Friend removed successfully",
	})
}

// GetPendingRequestsCount returns the count of pending friend requests
func GetPendingRequestsCount(c *fiber.Ctx) error {
	userID := GetUserId(c)
	if userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var count int64
	err := database.DB.Model(&models.Friendship{}).
		Where("friend_id = ? AND status = ?", userID, "pending").
		Count(&count).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get count",
		})
	}

	return c.JSON(fiber.Map{
		"count": count,
	})
}

// AreFriends checks if two users are friends
func AreFriends(userID1, userID2 uint) bool {
	// Same user is always considered "friend" (can see own comments)
	if userID1 == userID2 {
		return true
	}

	var friendship models.Friendship
	err := database.DB.Where(
		"((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?",
		userID1, userID2, userID2, userID1, "accepted",
	).First(&friendship).Error

	return err == nil
}

// getPendingRequestCount returns the count of pending friend requests for a user
func getPendingRequestCount(userID uint) int64 {
	var count int64
	database.DB.Model(&models.Friendship{}).
		Where("friend_id = ? AND status = ?", userID, "pending").
		Count(&count)
	return count
}

// getUsernameByID returns the username for a given user ID
func getUsernameByID(userID uint) string {
	var admin models.Admin
	database.DB.Select("username").First(&admin, userID)
	return admin.Username
}

// GetFriendIDs returns all friend IDs for a given user (including the user themselves)
func GetFriendIDs(userID uint) []uint {
	var friendships []models.Friendship
	database.DB.Where(
		"(user_id = ? OR friend_id = ?) AND status = ?",
		userID, userID, "accepted",
	).Find(&friendships)

	fmt.Printf("[DEBUG GetFriendIDs] userID=%d, found %d friendships\n", userID, len(friendships))

	friendIDs := []uint{userID} // Include self

	for _, fs := range friendships {
		fmt.Printf("[DEBUG GetFriendIDs] Friendship: ID=%d, UserID=%d, FriendID=%d, Status=%s\n", fs.ID, fs.UserID, fs.FriendID, fs.Status)
		if fs.UserID == userID {
			friendIDs = append(friendIDs, fs.FriendID)
		} else {
			friendIDs = append(friendIDs, fs.UserID)
		}
	}

	fmt.Printf("[DEBUG GetFriendIDs] Final friendIDs for user %d: %v\n", userID, friendIDs)
	return friendIDs
}
