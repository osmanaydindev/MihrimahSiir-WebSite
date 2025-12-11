package security

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// FileValidator provides file validation methods
type FileValidator struct{}

// NewFileValidator creates a new file validator instance
func NewFileValidator() *FileValidator {
	return &FileValidator{}
}

// ImageMagicBytes maps file extensions to their magic byte signatures
var ImageMagicBytes = map[string][]string{
	".jpg":  {"FFD8FF"},
	".jpeg": {"FFD8FF"},
	".png":  {"89504E47"},
	".gif":  {"474946383761", "474946383961"}, // GIF87a and GIF89a
	".webp": {"52494646"},                     // RIFF (WebP container)
}

// AllowedImageExtensions lists allowed image file extensions
var AllowedImageExtensions = []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

// ValidateImageFile validates an uploaded image file
func (fv *FileValidator) ValidateImageFile(file *multipart.FileHeader, maxSize int64) error {
	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", maxSize)
	}

	// Check file size is not zero
	if file.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	// Get file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext == "" {
		return fmt.Errorf("file has no extension")
	}

	// Check if extension is allowed
	isAllowedExt := false
	for _, allowedExt := range AllowedImageExtensions {
		if ext == allowedExt {
			isAllowedExt = true
			break
		}
	}

	if !isAllowedExt {
		return fmt.Errorf("file type not allowed. Allowed types: %s", strings.Join(AllowedImageExtensions, ", "))
	}

	// Open file to read magic bytes
	fileContent, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer fileContent.Close()

	// Read first 12 bytes (enough for all image formats)
	magicBytes := make([]byte, 12)
	n, err := fileContent.Read(magicBytes)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file content: %v", err)
	}

	// Convert to hex string
	hexStr := strings.ToUpper(hex.EncodeToString(magicBytes[:n]))

	// Validate magic bytes match the extension
	expectedMagicBytes, ok := ImageMagicBytes[ext]
	if !ok {
		return fmt.Errorf("unknown file extension: %s", ext)
	}

	// Check if any of the expected magic bytes match
	isValidMagic := false
	for _, expectedMagic := range expectedMagicBytes {
		if strings.HasPrefix(hexStr, expectedMagic) {
			isValidMagic = true
			break
		}
	}

	if !isValidMagic {
		return fmt.Errorf("file content does not match the expected file type (magic bytes mismatch)")
	}

	// Additional check for WebP
	if ext == ".webp" {
		// WebP should have "WEBP" after RIFF header
		if n >= 12 {
			webpSignature := string(magicBytes[8:12])
			if webpSignature != "WEBP" {
				return fmt.Errorf("invalid WebP file format")
			}
		}
	}

	// Check for dangerous content in file
	// Read more content to check for embedded scripts
	fileContent.Seek(0, 0) // Reset to beginning
	buffer := make([]byte, 1024)
	n, _ = fileContent.Read(buffer)

	if fv.containsDangerousImageContent(buffer[:n]) {
		return fmt.Errorf("file contains potentially dangerous content")
	}

	return nil
}

// containsDangerousImageContent checks for dangerous content in image files
func (fv *FileValidator) containsDangerousImageContent(content []byte) bool {
	contentStr := strings.ToLower(string(content))

	// Check for script tags
	dangerousPatterns := []string{
		"<script",
		"javascript:",
		"onerror=",
		"onload=",
		"<iframe",
		"<?php",
		"<?=",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(contentStr, pattern) {
			return true
		}
	}

	return false
}

// ValidateFilename validates and sanitizes a filename
func (fv *FileValidator) ValidateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename is empty")
	}

	// Check for path traversal attempts
	if strings.Contains(filename, "..") {
		return fmt.Errorf("filename contains path traversal sequence")
	}

	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("filename contains directory separator")
	}

	// Check for null bytes
	if bytes.ContainsRune([]byte(filename), '\x00') {
		return fmt.Errorf("filename contains null byte")
	}

	// Check length
	if len(filename) > 255 {
		return fmt.Errorf("filename too long (max 255 characters)")
	}

	return nil
}

// GetSafeExtension extracts and validates file extension
func (fv *FileValidator) GetSafeExtension(filename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filename))

	if ext == "" {
		return "", fmt.Errorf("no file extension found")
	}

	// Check if extension is in allowed list
	for _, allowedExt := range AllowedImageExtensions {
		if ext == allowedExt {
			return ext, nil
		}
	}

	return "", fmt.Errorf("file extension not allowed: %s", ext)
}
