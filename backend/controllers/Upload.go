package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"path/filepath"
)

type UploadService struct {
	BaseDir string
}

func NewUploadService(baseDir string) *UploadService {
	return &UploadService{BaseDir: baseDir}
}

func (s *UploadService) UploadFile(c *fiber.Ctx, fieldName string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", fmt.Errorf("dosya alınamadı: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	newFilename := uuid.New().String() + ext

	//fullDir := filepath.Join(s.BaseDir)
	//if err := os.MkdirAll(fullDir, os.ModePerm); err != nil {
	//	return "", fmt.Errorf("dizin oluşturulamadı: %w", err)
	//}

	fullPath := filepath.Join(s.BaseDir, "/", newFilename)
	if err := c.SaveFile(file, fullPath); err != nil {
		return "", fmt.Errorf("dosya kaydedilemedi: %w", err)
	}

	relativePath := filepath.Join(newFilename)
	return relativePath, nil
}
