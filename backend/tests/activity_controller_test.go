package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/models/activity"
)

func setupActivityTest() {
	testActivityModel.Clear()
	testActivityModel.SeedDefaultData()
}

func TestCreateActivityValid(t *testing.T) {
	setupActivityTest()
	validActivity := map[string]interface{}{
		"title":          "New Activity",
		"date":           "2025-12-31",
		"activity_image": "image_url",
		"description":    "Dpzinha legal demais",
	}

	body, _ := json.Marshal(validActivity)
	req, err := http.NewRequest("POST", "/activities", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestCreateActivityInvalid(t *testing.T) {
	setupActivityTest()
	invalidActivity := map[string]interface{}{
		"Dpzinha legal": "Title",
	}
	body, _ := json.Marshal(invalidActivity)
	req, err := http.NewRequest("POST", "/activities", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestReadActivity(t *testing.T) {
	setupActivityTest()
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

func TestUpdateActivityValid(t *testing.T) {
	setupActivityTest()
	validUpdate := map[string]interface{}{
		"description":    "Updated description",
		"activity_image": "https://wallsdesk.com/wp-content/uploads/2017/01/Octopus-Wallpapers-HD.jpg",
	}

	body, _ := json.Marshal(validUpdate)
	req, err := http.NewRequest("PUT", "/activities/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateActivityInvalid(t *testing.T) {
	setupActivityTest()
	invalidUpdate := map[string]interface{}{
		"name": "Invalid Update",
	}

	body, _ := json.Marshal(invalidUpdate)
	req, err := http.NewRequest("PUT", "/activities/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteActivityInvalid(t *testing.T) {
	setupActivityTest()
	invalidRequest := map[string]interface{}{
		"creator_id": 2,
	}

	body, _ := json.Marshal(invalidRequest)
	req, err := http.NewRequest("DELETE", "/activities/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}

func TestDeleteActivityValid(t *testing.T) {
	setupActivityTest()
	validRequest := map[string]interface{}{
		"creator_id": 1,
	}

	body, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("DELETE", "/activities/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testActivityRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
