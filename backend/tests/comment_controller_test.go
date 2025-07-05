package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend/models/comment"
)

func setupCommentTest() {
	setupActivityTest()
	testCommentModel.Clear()
	testCommentModel.SeedDefaultData()
}

func TestGetCommentsByActivityValid(t *testing.T) {
	setupCommentTest()
	req, err := http.NewRequest("GET", "/activities/1/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if activityID, ok := response["activity_id"].(float64); !ok || int(activityID) != 1 {
		t.Errorf("Expected activity_id to be 1, got %v", activityID)
	}

	if comments, ok := response["comments"].([]interface{}); !ok || len(comments) == 0 {
		t.Errorf("Expected comments array with at least one comment, got %v", comments)
	}
}

func TestGetCommentsByActivityNotFound(t *testing.T) {
	setupCommentTest()
	req, err := http.NewRequest("GET", "/activities/999/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestCreateCommentValid(t *testing.T) {
	setupCommentTest()
	validComment := map[string]interface{}{
		"user_id": 2,
		"content": "This is a test comment!",
	}

	body, _ := json.Marshal(validComment)
	req, err := http.NewRequest("POST", "/activities/1/comments", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response comment.Comment
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if response.UserID != 2 || response.Content != "This is a test comment!" || response.ActivityID != 1 {
		t.Errorf("Comment data mismatch: got %+v", response)
	}

	if response.ID == 0 || response.CreatedAt.IsZero() {
		t.Errorf("Comment should have ID and CreatedAt fields populated, %+v", response)
	}
}

func TestCreateCommentInvalidPayload(t *testing.T) {
	setupCommentTest()
	invalidComment := map[string]interface{}{
		"user_id": "2",
		// Missing content
	}

	body, _ := json.Marshal(invalidComment)
	req, err := http.NewRequest("POST", "/activities/1/comments", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestCreateCommentActivityNotFound(t *testing.T) {
	setupCommentTest()
	validComment := map[string]interface{}{
		"user_id": 2,
		"content": "This is a test comment!",
	}

	body, _ := json.Marshal(validComment)
	req, err := http.NewRequest("POST", "/activities/999/comments", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestDeleteCommentByAuthor(t *testing.T) {
	setupCommentTest()

	// First create a comment
	newComment := testCommentModel.CreateComment(comment.Comment{
		ActivityID: 1,
		UserID:     2,
		Content:    "Comment to be deleted",
	})

	deleteRequest := map[string]interface{}{
		"requester_id": 2, // Same user who created the comment
	}

	body, _ := json.Marshal(deleteRequest)
	req, err := http.NewRequest("DELETE", "/comments/"+strconv.Itoa(newComment.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify comment is deleted
	_, exists := testCommentModel.GetCommentByID(newComment.ID)
	if exists {
		t.Error("Comment should have been deleted")
	}
}

func TestDeleteCommentByActivityCreator(t *testing.T) {
	setupCommentTest()

	// Create a comment from user 2 on activity 1 (created by user 1)
	newComment := testCommentModel.CreateComment(comment.Comment{
		ActivityID: 1,
		UserID:     2,
		Content:    "Comment to be deleted by activity creator",
	})

	deleteRequest := map[string]interface{}{
		"requester_id": 1, // Activity creator (different from comment author)
	}

	body, _ := json.Marshal(deleteRequest)
	req, err := http.NewRequest("DELETE", "/comments/"+strconv.Itoa(newComment.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify comment is deleted
	_, exists := testCommentModel.GetCommentByID(newComment.ID)
	if exists {
		t.Error("Comment should have been deleted")
	}
}

func TestDeleteCommentForbidden(t *testing.T) {
	setupCommentTest()

	// Create a comment from user 2 on activity 1 (created by user 1)
	newComment := testCommentModel.CreateComment(comment.Comment{
		ActivityID: 1,
		UserID:     2,
		Content:    "Comment that should not be deletable by unauthorized user",
	})

	deleteRequest := map[string]interface{}{
		"requester_id": 3, // Different user (not comment author or activity creator)
	}

	body, _ := json.Marshal(deleteRequest)
	req, err := http.NewRequest("DELETE", "/comments/"+strconv.Itoa(newComment.ID), bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}

	// Verify comment still exists
	_, exists := testCommentModel.GetCommentByID(newComment.ID)
	if !exists {
		t.Error("Comment should not have been deleted")
	}
}

func TestDeleteCommentNotFound(t *testing.T) {
	setupCommentTest()

	deleteRequest := map[string]interface{}{
		"requester_id": 1,
	}

	body, _ := json.Marshal(deleteRequest)
	req, err := http.NewRequest("DELETE", "/comments/999", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testCommentRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}
