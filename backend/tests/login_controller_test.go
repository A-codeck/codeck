package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginSuccess(t *testing.T) {
	setupUserTest()
	loginRequest := map[string]interface{}{
		"email":    "user@example.com",
		"password": "password123",
	}

	body, _ := json.Marshal(loginRequest)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testLoginRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if _, exists := response["token"]; !exists {
		t.Error("No token in login response")
	}

	if user, exists := response["user"].(map[string]interface{}); !exists {
		t.Error("No user info in login response")
	} else {
		if user["id"] == nil || user["name"] == nil || user["email"] == nil {
			t.Error("Missing required user fields in login response")
		}
		if _, exists := user["password"]; exists {
			t.Error("Password should not be returned in login response")
		}
	}
}

func TestLoginInvalidCredentials(t *testing.T) {
	setupUserTest()
	loginRequest := map[string]interface{}{
		"email":    "user@example.com",
		"password": "wrongpassword",
	}

	body, _ := json.Marshal(loginRequest)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testLoginRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

func TestLoginMissingFields(t *testing.T) {
	setupUserTest()
	loginRequest := map[string]interface{}{
		"email": "user@example.com",
		// Missing password
	}

	body, _ := json.Marshal(loginRequest)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testLoginRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
