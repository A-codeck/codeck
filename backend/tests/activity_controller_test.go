package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"backend/models/activity"
	"backend/utils"
)

func setupActivityTest(t *testing.T) {
	// Use comprehensive seeding that includes all dependencies
	SeedAllTestData(t)

	// Initialize upload directory for tests
	if err := utils.InitUploadDir(); err != nil {
		t.Fatalf("Failed to initialize upload directory: %v", err)
	}
}

// createTestImageFile creates a temporary image file for testing
func createTestImageFile() (*os.File, error) {
	// Create a simple 1x1 PNG file for testing
	pngData := []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, // IHDR chunk
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, // 1x1 dimensions
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53, 0xDE, // rest of IHDR
		0x00, 0x00, 0x00, 0x0C, 0x49, 0x44, 0x41, 0x54, // IDAT chunk
		0x08, 0x99, 0x01, 0x01, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82, // IEND chunk
	}

	tmpfile, err := os.CreateTemp("", "test-image-*.png")
	if err != nil {
		return nil, err
	}

	if _, err := tmpfile.Write(pngData); err != nil {
		tmpfile.Close()
		os.Remove(tmpfile.Name())
		return nil, err
	}

	// Seek back to beginning for reading
	tmpfile.Seek(0, 0)
	return tmpfile, nil
}

// createMultipartFormData creates a multipart form with the given fields and file
func createMultipartFormData(fields map[string]string, imageFile *os.File) (*bytes.Buffer, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return nil, "", err
		}
	}

	// Add image file
	if imageFile != nil {
		// Create form file with correct content type
		part, err := writer.CreateFormFile("image", filepath.Base(imageFile.Name()))
		if err != nil {
			return nil, "", err
		}

		if _, err := io.Copy(part, imageFile); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &buf, writer.FormDataContentType(), nil
}

func TestCreateActivityValid(t *testing.T) {
	setupActivityTest(t)

	// Create test image file
	imageFile, err := createTestImageFile()
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer os.Remove(imageFile.Name())
	defer imageFile.Close()

	// Prepare form data
	fields := map[string]string{
		"creator_id":  "1",
		"group_id":    "1",
		"title":       "New Activity",
		"date":        "2025-12-31",
		"description": "Dpzinha legal demais",
	}

	body, contentType, err := createMultipartFormData(fields, imageFile)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req, err := http.NewRequest("POST", "/activities", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v. Response: %s",
			status, http.StatusCreated, recorder.Body.String())
	}

	// Verify response contains the activity
	var activity activity.Activity
	if err := json.Unmarshal(recorder.Body.Bytes(), &activity); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if activity.Title != "New Activity" {
		t.Errorf("Expected title 'New Activity', got '%s'", activity.Title)
	}

	if activity.ActivityImage == "" {
		t.Error("Expected activity_image to be set")
	}

	// Cleanup uploaded file
	if activity.ActivityImage != "" {
		filename := strings.TrimPrefix(activity.ActivityImage, "/uploads/activities/")
		utils.DeleteImage(filename)
	}
}

func TestCreatePersonalActivityValid(t *testing.T) {
	setupActivityTest(t)

	// Create test image file
	imageFile, err := createTestImageFile()
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer os.Remove(imageFile.Name())
	defer imageFile.Close()

	// Prepare form data for personal activity (group_id = 0)
	fields := map[string]string{
		"creator_id":  "1",
		"group_id":    "0",
		"title":       "Personal Activity",
		"date":        "2025-12-31",
		"description": "Personal coding session",
	}

	body, contentType, err := createMultipartFormData(fields, imageFile)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req, err := http.NewRequest("POST", "/activities", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v. Response: %s",
			status, http.StatusCreated, recorder.Body.String())
	}

	// Verify response
	var activity activity.Activity
	if err := json.Unmarshal(recorder.Body.Bytes(), &activity); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if activity.GroupID != 0 {
		t.Errorf("Expected group_id 0 for personal activity, got %d", activity.GroupID)
	}

	// Cleanup uploaded file
	if activity.ActivityImage != "" {
		filename := strings.TrimPrefix(activity.ActivityImage, "/uploads/activities/")
		utils.DeleteImage(filename)
	}
}

func TestCreateActivityWithoutImage(t *testing.T) {
	setupActivityTest(t)

	// Prepare form data without image
	fields := map[string]string{
		"creator_id":  "1",
		"group_id":    "1",
		"title":       "Activity Without Image",
		"date":        "2025-12-31",
		"description": "This should fail",
	}

	body, contentType, err := createMultipartFormData(fields, nil)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req, err := http.NewRequest("POST", "/activities", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v. Expected failure due to missing image.",
			status, http.StatusBadRequest)
	}
}

func TestCreateActivityInvalidFields(t *testing.T) {
	setupActivityTest(t)

	// Create test image file
	imageFile, err := createTestImageFile()
	if err != nil {
		t.Fatalf("Failed to create test image: %v", err)
	}
	defer os.Remove(imageFile.Name())
	defer imageFile.Close()

	// Prepare form data with missing required fields
	fields := map[string]string{
		"creator_id": "1",
		// Missing title and date
		"description": "Missing required fields",
	}

	body, contentType, err := createMultipartFormData(fields, imageFile)
	if err != nil {
		t.Fatalf("Failed to create multipart form: %v", err)
	}

	req, err := http.NewRequest("POST", "/activities", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestReadActivity(t *testing.T) {
	setupActivityTest(t)
	req, err := http.NewRequest("GET", "/activities/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var activity activity.Activity
	if err := json.NewDecoder(recorder.Body).Decode(&activity); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if activity.ID == 0 || activity.Title == "" || activity.Date.IsZero() {
		t.Error("Missing required activity fields in response")
	}
}

// Comprehensive tests for activity feed functionality
func TestGetUserFeedHasAllActivities(t *testing.T) {
	setupActivityTest(t)

	// Now get the feed for user 1
	feedReq, err := http.NewRequest("GET", "/activities/feed?user_id=1", nil)
	if err != nil {
		t.Fatal(err)
	}
	feedRecorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(feedRecorder, feedReq)
	if status := feedRecorder.Code; status != http.StatusOK {
		t.Fatalf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var activities []activity.Activity
	if err := json.NewDecoder(feedRecorder.Body).Decode(&activities); err != nil {
		t.Fatal("Failed to decode feed response body")
	}

	// Should have at least the seeded activities
	if len(activities) == 0 {
		t.Error("Expected at least some activities in user feed")
	}
}
