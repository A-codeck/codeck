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
		// Missing data
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
		// Missing required fields
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

func TestAddUserToGroupValid(t *testing.T) {
	setupGroupTest()
	validRequest := map[string]interface{}{
		"user_id":  "2",
		"nickname": "TestUser",
	}

	body, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
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

func TestAddUserToGroupInvalid(t *testing.T) {
	setupGroupTest()
	invalidRequest := map[string]interface{}{
		"random": "2",
	}

	body, _ := json.Marshal(invalidRequest)
	req, err := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
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

func TestAddUserToGroupDuplicate(t *testing.T) {
	setupGroupTest()
	// Add user first time
	validRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(validRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Try to add same user again
	body, _ = json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}
}

func TestRemoveUserFromGroupValid(t *testing.T) {
	setupGroupTest()
	// First add a user
	addRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Now remove the user
	removeRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "1", // Group creator
	}
	body, _ = json.Marshal(removeRequest)
	req, err := http.NewRequest("DELETE", "/groups/1/members", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestRemoveUserFromGroupForbidden(t *testing.T) {
	setupGroupTest()
	// First add a user
	addRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Try to remove user with unauthorized requester
	removeRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "3", // Not group creator or the user themselves
	}
	body, _ = json.Marshal(removeRequest)
	req, err := http.NewRequest("DELETE", "/groups/1/members", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}

func TestGetGroupMembers(t *testing.T) {
	setupGroupTest()
	// Add some users to the group
	users := []string{"2", "3"}
	for _, userID := range users {
		addRequest := map[string]interface{}{
			"user_id": userID,
		}
		body, _ := json.Marshal(addRequest)
		req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		testGroupRouter.ServeHTTP(recorder, req)
	}

	// Get group members
	req, err := http.NewRequest("GET", "/groups/1/members", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if memberCount, ok := response["member_count"].(float64); !ok || memberCount != 3 {
		t.Errorf("Expected member_count to be 3 (creator + 2 added), got %v", memberCount)
	}
}

func TestCreateInviteLinkValid(t *testing.T) {
	setupGroupTest()
	validRequest := map[string]interface{}{
		"creator_id": "1", // Group creator
	}

	body, _ := json.Marshal(validRequest)
	req, err := http.NewRequest("POST", "/groups/1/invites", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var invite group.GroupInvite
	if err := json.NewDecoder(recorder.Body).Decode(&invite); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if invite.InviteCode == "" || invite.GroupID != "1" || invite.CreatedBy != "1" {
		t.Error("Invalid invite data in response")
	}
}

func TestCreateInviteLinkForbidden(t *testing.T) {
	setupGroupTest()
	invalidRequest := map[string]interface{}{
		"creator_id": "2", // Not group creator
	}

	body, _ := json.Marshal(invalidRequest)
	req, err := http.NewRequest("POST", "/groups/1/invites", bytes.NewBuffer(body))
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

func TestJoinGroupByInviteValid(t *testing.T) {
	setupGroupTest()

	// First create an invite
	createRequest := map[string]interface{}{
		"creator_id": "1",
	}
	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/groups/1/invites", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	var invite group.GroupInvite
	json.NewDecoder(recorder.Body).Decode(&invite)

	// Now use the invite to join
	joinRequest := map[string]interface{}{
		"user_id":  "2",
		"nickname": "NewMember",
	}
	body, _ = json.Marshal(joinRequest)
	req, err := http.NewRequest("POST", "/invites/"+invite.InviteCode+"/join", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestJoinGroupByInviteInvalidCode(t *testing.T) {
	setupGroupTest()

	joinRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(joinRequest)
	req, err := http.NewRequest("POST", "/invites/INVALID123/join", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestJoinGroupByInviteDuplicateUser(t *testing.T) {
	setupGroupTest()

	// First create an invite
	createRequest := map[string]interface{}{
		"creator_id": "1",
	}
	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/groups/1/invites", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	var invite group.GroupInvite
	json.NewDecoder(recorder.Body).Decode(&invite)

	// Join first time
	joinRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ = json.Marshal(joinRequest)
	req, _ = http.NewRequest("POST", "/invites/"+invite.InviteCode+"/join", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Try to join again with same user
	body, _ = json.Marshal(joinRequest)
	req, err := http.NewRequest("POST", "/invites/"+invite.InviteCode+"/join", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}
}

func TestGetGroupInvites(t *testing.T) {
	setupGroupTest()

	// Create a couple of invites
	for i := 0; i < 2; i++ {
		createRequest := map[string]interface{}{
			"creator_id": "1",
		}
		body, _ := json.Marshal(createRequest)
		req, _ := http.NewRequest("POST", "/groups/1/invites", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		testGroupRouter.ServeHTTP(recorder, req)
	}

	// Get invites
	req, err := http.NewRequest("GET", "/groups/1/invites", nil)
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if inviteCount, ok := response["invite_count"].(float64); !ok || inviteCount != 2 {
		t.Errorf("Expected invite_count to be 2, got %v", inviteCount)
	}
}

func TestDeactivateInviteValid(t *testing.T) {
	setupGroupTest()

	// First create an invite
	createRequest := map[string]interface{}{
		"creator_id": "1",
	}
	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/groups/1/invites", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	var invite group.GroupInvite
	json.NewDecoder(recorder.Body).Decode(&invite)

	// Deactivate the invite
	deactivateRequest := map[string]interface{}{
		"requester_id": "1", // Group creator
	}
	body, _ = json.Marshal(deactivateRequest)
	req, err := http.NewRequest("DELETE", "/invites/"+invite.InviteCode+"/deactivate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestDeactivateInviteForbidden(t *testing.T) {
	setupGroupTest()

	// First create an invite
	createRequest := map[string]interface{}{
		"creator_id": "1",
	}
	body, _ := json.Marshal(createRequest)
	req, _ := http.NewRequest("POST", "/groups/1/invites", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	var invite group.GroupInvite
	json.NewDecoder(recorder.Body).Decode(&invite)

	// Try to deactivate with unauthorized user
	deactivateRequest := map[string]interface{}{
		"requester_id": "3", // Not group creator or invite creator
	}
	body, _ = json.Marshal(deactivateRequest)
	req, err := http.NewRequest("DELETE", "/invites/"+invite.InviteCode+"/deactivate", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}

func TestSetUserNicknameValid(t *testing.T) {
	setupGroupTest()
	// First add a user to the group
	addRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Now set the nickname
	nicknameRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "2", // User themselves
		"nickname":     "CoolNickname",
	}
	body, _ = json.Marshal(nicknameRequest)
	req, err := http.NewRequest("PUT", "/groups/1/members/nickname", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if nickname, ok := response["nickname"].(string); !ok || nickname != "CoolNickname" {
		t.Errorf("Expected nickname to be 'CoolNickname', got %v", nickname)
	}
}

func TestSetUserNicknameByGroupCreator(t *testing.T) {
	setupGroupTest()
	// First add a user to the group
	addRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Now set the nickname as group creator
	nicknameRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "1", // Group creator
		"nickname":     "AssignedByOwner",
	}
	body, _ = json.Marshal(nicknameRequest)
	req, err := http.NewRequest("PUT", "/groups/1/members/nickname", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestSetUserNicknameForbidden(t *testing.T) {
	setupGroupTest()
	// First add a user to the group
	addRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Try to set nickname as unauthorized user
	nicknameRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "3", // Neither the user nor group creator
		"nickname":     "Unauthorized",
	}
	body, _ = json.Marshal(nicknameRequest)
	req, err := http.NewRequest("PUT", "/groups/1/members/nickname", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}

func TestSetUserNicknameUserNotInGroup(t *testing.T) {
	setupGroupTest()

	// Try to set nickname for a user not in the group
	nicknameRequest := map[string]interface{}{
		"user_id":      "3", // User not in group
		"requester_id": "1", // Group creator
		"nickname":     "NotInGroup",
	}
	body, _ := json.Marshal(nicknameRequest)
	req, err := http.NewRequest("PUT", "/groups/1/members/nickname", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestSetUserNicknameTooLong(t *testing.T) {
	setupGroupTest()
	// First add a user to the group
	addRequest := map[string]interface{}{
		"user_id": "2",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Try to set a nickname that's too long
	longNickname := make([]byte, 51) // 51 characters
	for i := range longNickname {
		longNickname[i] = 'a'
	}

	nicknameRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "2",
		"nickname":     string(longNickname),
	}
	body, _ = json.Marshal(nicknameRequest)
	req, err := http.NewRequest("PUT", "/groups/1/members/nickname", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestDeleteUserNicknameValid(t *testing.T) {
	setupGroupTest()
	// First add a user with a nickname
	addRequest := map[string]interface{}{
		"user_id":  "2",
		"nickname": "InitialNickname",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Now delete the nickname
	deleteRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "2", // User themselves
	}
	body, _ = json.Marshal(deleteRequest)
	req, err := http.NewRequest("DELETE", "/groups/1/members/nickname", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response body")
	}

	if nickname := response["nickname"]; nickname != nil {
		t.Errorf("Expected nickname to be nil after deletion, got %v", nickname)
	}
}

func TestDeleteUserNicknameForbidden(t *testing.T) {
	setupGroupTest()
	// First add a user with a nickname
	addRequest := map[string]interface{}{
		"user_id":  "2",
		"nickname": "InitialNickname",
	}
	body, _ := json.Marshal(addRequest)
	req, _ := http.NewRequest("POST", "/groups/1/members", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	// Try to delete nickname as unauthorized user
	deleteRequest := map[string]interface{}{
		"user_id":      "2",
		"requester_id": "3", // Neither the user nor group creator
	}
	body, _ = json.Marshal(deleteRequest)
	req, err := http.NewRequest("DELETE", "/groups/1/members/nickname", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	testGroupRouter.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusForbidden)
	}
}
