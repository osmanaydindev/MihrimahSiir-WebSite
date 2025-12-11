package controllers

import (
	"backend/database"
	"backend/models"
	"github.com/gofiber/fiber/v2"
	"time"
)

func CreateLog(c *fiber.Ctx, username string, roleId uint, userId uint) error {

	nowDate := time.Now().Format("2006-01-02 15:04:05")
	log := models.Log{
		Username:  username,
		RoleID:    roleId,
		UserId:    userId,
		LoginDate: nowDate,
	}
	err := database.DB.Create(&log)
	if err.Error != nil {
		println(err.Error)
	}
	if log.ID == 0 {
		return c.SendStatus(400)
	}
	return c.JSON(log)
}
func GetLogs(c *fiber.Ctx) error {
	var logs []models.Log

	database.DB.Find(&logs)
	return c.JSON(logs)
}
func DeleteLog(c *fiber.Ctx) error {
	id := c.Params("id")
	var log models.Log
	database.DB.First(&log, id)
	if log.ID == 0 {
		return c.Status(404).SendString("Log not found")
	}
	database.DB.Delete(&log)

	var logs []models.Log
	database.DB.Find(&logs)
	return c.JSON(logs)
}
