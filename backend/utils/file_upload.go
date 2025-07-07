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

const (
	MaxFileSize = 10 << 20 // 10MB
	UploadDir   = "uploads/activities"
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
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
}

// InitUploadDir creates the upload directory if it doesn't exist
func InitUploadDir() error {
	return os.MkdirAll(UploadDir, 0755)
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

// SaveImage saves the uploaded image file to the server
func SaveImage(file multipart.File, header *multipart.FileHeader) (*FileUploadResult, error) {
	// Validate the image
	if err := ValidateImage(file, header); err != nil {
		return nil, err
	}

	// Reset file pointer to beginning
	file.Seek(0, 0)

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	uniqueID := uuid.New().String()
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%s_%d%s", uniqueID, timestamp, ext)

	// Create full path
	fullPath := filepath.Join(UploadDir, filename)

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
	}, nil
}

// DeleteImage removes an image file from the server
func DeleteImage(filename string) error {
	if filename == "" {
		return nil
	}

	// Extract just the filename from URL if needed
	if strings.HasPrefix(filename, "/") {
		filename = filepath.Base(filename)
	}

	fullPath := filepath.Join(UploadDir, filename)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	return os.Remove(fullPath)
}
