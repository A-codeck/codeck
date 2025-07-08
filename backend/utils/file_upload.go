package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ImageType represents the type of image being uploaded
type ImageType string

const (
	// Image types
	ActivityImageType ImageType = "activity"
	GroupImageType    ImageType = "group"

	// File constraints
	MaxFileSize = 10 << 20 // 10MB

	// Upload directories
	ActivityUploadDir = "uploads/activities"
	GroupUploadDir    = "uploads/groups"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// FileUploadResult represents the result of a file upload
type FileUploadResult struct {
	Filename     string    `json:"filename"`
	OriginalName string    `json:"original_name"`
	Size         int64     `json:"size"`
	URL          string    `json:"url"`
	Type         ImageType `json:"type"`
}

// GetUploadDir returns the upload directory for the given image type
func GetUploadDir(imageType ImageType) string {
	switch imageType {
	case ActivityImageType:
		return ActivityUploadDir
	case GroupImageType:
		return GroupUploadDir
	default:
		return ActivityUploadDir // Default fallback
	}
}

// InitUploadDirs creates the upload directories if they don't exist
func InitUploadDirs() error {
	if err := os.MkdirAll(ActivityUploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create activity upload directory: %v", err)
	}
	if err := os.MkdirAll(GroupUploadDir, 0755); err != nil {
		return fmt.Errorf("failed to create group upload directory: %v", err)
	}
	return nil
}

// ValidateImage validates if the uploaded file is a valid image
func ValidateImage(file multipart.File, header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", MaxFileSize)
	}

	// Check file extension (primary validation)
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	if !validExtensions[ext] {
		return fmt.Errorf("invalid file extension. Allowed extensions: .jpg, .jpeg, .png, .gif, .webp")
	}

	return nil
}

// generateUniqueFilename creates a unique filename with timestamp and UUID
func generateUniqueFilename(originalFilename string, imageType ImageType) string {
	ext := filepath.Ext(originalFilename)
	uniqueID := uuid.New().String()
	timestamp := time.Now().Unix()

	// Add type prefix for better organization
	var prefix string
	switch imageType {
	case ActivityImageType:
		prefix = "activity"
	case GroupImageType:
		prefix = "group"
	default:
		prefix = "unknown"
	}

	return fmt.Sprintf("%s_%s_%d%s", prefix, uniqueID, timestamp, ext)
}

// SaveImage saves the uploaded image file to the server with the specified type
func SaveImage(file multipart.File, header *multipart.FileHeader, imageType ImageType) (*FileUploadResult, error) {
	// Validate the image
	if err := ValidateImage(file, header); err != nil {
		return nil, err
	}

	// Reset file pointer to beginning
	file.Seek(0, 0)

	// Generate unique filename with type prefix
	filename := generateUniqueFilename(header.Filename, imageType)
	uploadDir := GetUploadDir(imageType)
	fullPath := filepath.Join(uploadDir, filename)

	// Create destination file
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %v", err)
	}
	defer dst.Close()

	// Copy file content
	size, err := io.Copy(dst, file)
	if err != nil {
		// Clean up created file on error
		os.Remove(fullPath)
		return nil, fmt.Errorf("failed to save file: %v", err)
	}

	// Generate URL (relative path that can be served by static file handler)
	url := fmt.Sprintf("/%s", fullPath)

	return &FileUploadResult{
		Filename:     filename,
		OriginalName: header.Filename,
		Size:         size,
		URL:          url,
		Type:         imageType,
	}, nil
}

// SaveActivityImage saves an activity image (convenience function)
func SaveActivityImage(file multipart.File, header *multipart.FileHeader) (*FileUploadResult, error) {
	return SaveImage(file, header, ActivityImageType)
}

// SaveGroupImage saves a group image (convenience function)
func SaveGroupImage(file multipart.File, header *multipart.FileHeader) (*FileUploadResult, error) {
	return SaveImage(file, header, GroupImageType)
}

// DeleteImage removes an image file from the server
func DeleteImage(filename string, imageType ImageType) error {
	if filename == "" {
		return nil
	}

	// Extract just the filename from URL if needed
	if strings.HasPrefix(filename, "/") {
		filename = filepath.Base(filename)
	}

	uploadDir := GetUploadDir(imageType)
	fullPath := filepath.Join(uploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	return os.Remove(fullPath)
}

// DeleteActivityImage removes an activity image (convenience function)
func DeleteActivityImage(filename string) error {
	return DeleteImage(filename, ActivityImageType)
}

// DeleteGroupImage removes a group image (convenience function)
func DeleteGroupImage(filename string) error {
	return DeleteImage(filename, GroupImageType)
}
