package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/models/activity"
)

func setupUserTest() {
	testUserModel.Clear()
	testUserModel.SeedDefaultData()
}

func TestCreateUserValid(t *testing.T) {
	setupUserTest()
	userData := map[string]interface{}{
		"email":    "newuser@example.com",
		"name":     "New User",
		"password": "securepassword123",
	}

	body, _ := json.Marshal(userData)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testUserRouter.ServeHTTP(recorder, req)

	t.Logf("Response status: %d", recorder.Code)
	t.Logf("Response body: %s", recorder.Body.String())

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdUser map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&createdUser); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	if createdUser["id"] == nil || createdUser["email"] == nil || createdUser["name"] == nil {
		t.Error("Missing required user fields in response")
	}

	if _, exists := createdUser["password"]; exists {
		t.Error("Password should not be returned in response")
	}
}

func TestCreateUserInvalid(t *testing.T) {
	setupUserTest()

	invalidUser := map[string]interface{}{
		"name": "Just a name",
	}
	body, _ := json.Marshal(invalidUser)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testUserRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestGetUserInfo(t *testing.T) {
	setupUserTest()
	req, err := http.NewRequest("GET", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testUserRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var userData map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&userData); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	if userData["id"] == nil || userData["name"] == nil || userData["email"] == nil {
		t.Error("Missing required user fields in response")
	}

	if _, exists := userData["password"]; exists {
		t.Error("Password should not be returned in response")
	}
}

func TestGetUserActivities(t *testing.T) {
	setupUserTest()
	testActivityModel.Clear()
	testActivityModel.SeedDefaultData()

	testActivityModel.CreateActivity(activity.Activity{
		CreatorID:   "1",
		Title:       "Test Activity 1",
		Date:        "2025-07-01",
		Description: stringPtr("Test description 1"),
	})

	testActivityModel.CreateActivity(activity.Activity{
		CreatorID:   "1",
		Title:       "Test Activity 2",
		Date:        "2025-07-02",
		Description: stringPtr("Test description 2"),
	})

	testActivityModel.CreateActivity(activity.Activity{
		CreatorID:   "999",
		Title:       "Other User Activity",
		Date:        "2025-07-03",
		Description: stringPtr("Should not be returned"),
	})

	req, err := http.NewRequest("GET", "/users/1/activities", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testUserRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var activities []map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&activities); err != nil {
		t.Fatal("Failed to decode response body:", err)
	}

	// Should have 3 activities for user "1" (1 from seed + 2 created in test)
	if len(activities) != 3 {
		t.Errorf("Expected 3 activities, got %d", len(activities))
	}

	// Verify all activities belong to user "1"
	for _, activity := range activities {
		if creatorID, ok := activity["creator_id"].(string); !ok || creatorID != "1" {
			t.Errorf("Activity creator_id should be '1', got %v", activity["creator_id"])
		}
	}
}

func TestGetUserActivitiesUserNotFound(t *testing.T) {
	setupUserTest()

	req, err := http.NewRequest("GET", "/users/999/activities", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testUserRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func stringPtr(s string) *string {
	return &s
}
