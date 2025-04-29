package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/models/group"
)

func setupGroupTest() {
	testGroupModel.Clear()
	testGroupModel.SeedDefaultData()
}

func TestCreateGroupValid(t *testing.T) {
	setupGroupTest()
	validGroup := map[string]interface{}{
		"name":        "New Group",
		"end_date":    "2025-12-31",
		"group_image": "image_url",
		"description": "A test group",
	}

	body, _ := json.Marshal(validGroup)
	req, err := http.NewRequest("POST", "/groups", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestCreateGroupInvalid(t *testing.T) {
	setupGroupTest()
	invalidGroup := map[string]interface{}{
		"end_date": "2025-12-31",
	}
	body, _ := json.Marshal(invalidGroup)
	req, err := http.NewRequest("POST", "/groups", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestReadGroup(t *testing.T) {
	setupGroupTest()
	req, err := http.NewRequest("GET", "/groups/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var group group.Group
	if err := json.NewDecoder(recorder.Body).Decode(&group); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if group.ID == "" || group.Name == "" || group.StartDate == "" || group.EndDate == "" {
		t.Error("Missing required group fields in response")
	}
}

func TestUpdateGroupValid(t *testing.T) {
	setupGroupTest()
	validUpdate := map[string]interface{}{
		"description": "Updated description",
		"group_image": "new_image_url",
		"end_date":    "2026-01-01",
	}

	body, _ := json.Marshal(validUpdate)
	req, err := http.NewRequest("PUT", "/groups/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestUpdateGroupInvalid(t *testing.T) {
	setupGroupTest()
	invalidUpdate := map[string]interface{}{
		"name": "Invalid Update",
	}

	body, _ := json.Marshal(invalidUpdate)
	req, err := http.NewRequest("PUT", "/groups/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteGroupInvalid(t *testing.T) {
	setupGroupTest()
	invalidRequest := map[string]interface{}{
		"creator_id": "2",
	}

	body, _ := json.Marshal(invalidRequest)
	req, err := http.NewRequest("DELETE", "/groups/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}

func TestDeleteGroupValid(t *testing.T) {
	setupGroupTest()
	validRequest := map[string]interface{}{
		"creator_id": "1",
	}

	body, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("DELETE", "/groups/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}
