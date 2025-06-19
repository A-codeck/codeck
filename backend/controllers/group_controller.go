package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/models/group"

	"github.com/gorilla/mux"
)

type GroupController struct {
	Model group.GroupModel
}

func NewGroupController(model group.GroupModel) *GroupController {
	return &GroupController{Model: model}
}

func (gc *GroupController) GetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]
	
	requesterID := r.URL.Query().Get("requester_id")
	if requesterID == "" {
		http.Error(w, "Missing requester_id", http.StatusBadRequest)
		return
	}
	
	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, requesterID) {
		http.Error(w, "Forbidden: Only group members can view group details", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}

func (gc *GroupController) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var group group.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if group.Name == "" || group.EndDate == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	createdGroup := gc.Model.CreateGroup(group)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGroup)
}

func (gc *GroupController) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, ok := updates["name"].(string); ok {
		http.Error(w, "Name field cannot be updated", http.StatusBadRequest)
		return
	}

	updatedGroup, exists := gc.Model.UpdateGroup(groupID, updates)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedGroup)
}

func (gc *GroupController) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]
	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Failed to decode request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	creatorID, ok := request["creator_id"].(string)
	if !ok {
		http.Error(w, "Missing creator_id", http.StatusBadRequest)
		return
	}

	if creatorID != group.CreatorID {
		// Group exists but creator_id is invalid
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if !gc.Model.DeleteGroup(groupID) {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (gc *GroupController) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	_, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	if gc.Model.IsUserInGroup(groupID, request.UserID) {
		http.Error(w, "User is already a member of this group", http.StatusConflict)
		return
	}

	success := gc.Model.AddUserToGroup(groupID, request.UserID)
	if !success {
		http.Error(w, "Failed to add user to group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := map[string]interface{}{
		"message":  "User added to group successfully",
		"group_id": groupID,
		"user_id":  request.UserID,
	}

	json.NewEncoder(w).Encode(response)
}

func (gc *GroupController) RemoveUserFromGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID      string `json:"user_id"`
		RequesterID string `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == "" || request.RequesterID == "" {
		http.Error(w, "Missing user_id or requester_id", http.StatusBadRequest)
		return
	}

	if request.RequesterID != group.CreatorID && request.RequesterID != request.UserID {
		http.Error(w, "Forbidden: Only group creator or the user themselves can remove membership", http.StatusForbidden)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, request.UserID) {
		http.Error(w, "User is not a member of this group", http.StatusNotFound)
		return
	}

	success := gc.Model.RemoveUserFromGroup(groupID, request.UserID)
	if !success {
		http.Error(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":  "User removed from group successfully",
		"group_id": groupID,
		"user_id":  request.UserID,
	})
}

func (gc *GroupController) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]
	
	requesterID := r.URL.Query().Get("requester_id")
	if requesterID == "" {
		http.Error(w, "Missing requester_id", http.StatusBadRequest)
		return
	}

	_, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, requesterID) {
		http.Error(w, "Forbidden: Only group members can view group members", http.StatusForbidden)
		return
	}

	members, exists := gc.Model.GetGroupMembers(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"group_id":     groupID,
		"members":      members,
		"member_count": len(members),
	})
}

func (gc *GroupController) CreateInviteLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		CreatorID string  `json:"creator_id"`
		ExpiresAt *string `json:"expires_at,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.CreatorID == "" {
		http.Error(w, "Missing creator_id", http.StatusBadRequest)
		return
	}

	if request.CreatorID != group.CreatorID {
		http.Error(w, "Forbidden: Only group creator can create invite links", http.StatusForbidden)
		return
	}

	invite, success := gc.Model.CreateInviteLink(groupID, request.CreatorID, request.ExpiresAt)
	if !success {
		http.Error(w, "Failed to create invite link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(invite)
}

func (gc *GroupController) JoinGroupByInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inviteCode := vars["invite_code"]

	var request struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	invite, exists := gc.Model.GetInviteByCode(inviteCode)
	if !exists || !invite.IsActive {
		http.Error(w, "Invalid or expired invite code", http.StatusNotFound)
		return
	}

	if gc.Model.IsUserInGroup(invite.GroupID, request.UserID) {
		http.Error(w, "User is already a member of this group", http.StatusConflict)
		return
	}

	success := gc.Model.AddUserToGroup(invite.GroupID, request.UserID)
	if !success {
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":     "Successfully joined group",
		"group_id":    invite.GroupID,
		"user_id":     request.UserID,
		"invite_code": inviteCode,
	}
	
	json.NewEncoder(w).Encode(response)
}

func (gc *GroupController) GetGroupInvites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	_, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	invites := gc.Model.GetActiveInvites(groupID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"group_id":     groupID,
		"invites":      invites,
		"invite_count": len(invites),
	})
}

func (gc *GroupController) DeactivateInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inviteCode := vars["invite_code"]

	var request struct {
		RequesterID string `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.RequesterID == "" {
		http.Error(w, "Missing requester_id", http.StatusBadRequest)
		return
	}

	invite, exists := gc.Model.GetInviteByCode(inviteCode)
	if !exists {
		http.Error(w, "Invite not found", http.StatusNotFound)
		return
	}

	group, exists := gc.Model.GetGroupByID(invite.GroupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if request.RequesterID != group.CreatorID && request.RequesterID != invite.CreatedBy {
		http.Error(w, "Forbidden: Only group creator or invite creator can deactivate invite", http.StatusForbidden)
		return
	}

	success := gc.Model.DeactivateInvite(inviteCode)
	if !success {
		http.Error(w, "Failed to deactivate invite", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":     "Invite deactivated successfully",
		"invite_code": inviteCode,
	})
}

func (gc *GroupController) SetUserNickname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID      string  `json:"user_id"`
		RequesterID string  `json:"requester_id"`
		Nickname    *string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == "" || request.RequesterID == "" {
		http.Error(w, "Missing user_id or requester_id", http.StatusBadRequest)
		return
	}

	if request.RequesterID != request.UserID && request.RequesterID != group.CreatorID {
		http.Error(w, "Forbidden: Only the user themselves or group creator can set nickname", http.StatusForbidden)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, request.UserID) {
		http.Error(w, "User is not a member of this group", http.StatusNotFound)
		return
	}

	if request.Nickname != nil && len(*request.Nickname) == 0 {
		request.Nickname = nil
	}
	if request.Nickname != nil && len(*request.Nickname) > 50 {
		http.Error(w, "Nickname cannot be longer than 50 characters", http.StatusBadRequest)
		return
	}

	success := gc.Model.SetUserNickname(groupID, request.UserID, request.Nickname)
	if !success {
		http.Error(w, "Failed to set nickname", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":  "Nickname updated successfully",
		"group_id": groupID,
		"user_id":  request.UserID,
	}
	if request.Nickname != nil {
		response["nickname"] = *request.Nickname
	} else {
		response["nickname"] = nil
	}
	json.NewEncoder(w).Encode(response)
}

func (gc *GroupController) DeleteUserNickname(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID      string `json:"user_id"`
		RequesterID string `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == "" || request.RequesterID == "" {
		http.Error(w, "Missing user_id or requester_id", http.StatusBadRequest)
		return
	}

	if request.RequesterID != request.UserID && request.RequesterID != group.CreatorID {
		http.Error(w, "Forbidden: Only the user themselves or group creator can delete nickname", http.StatusForbidden)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, request.UserID) {
		http.Error(w, "User is not a member of this group", http.StatusNotFound)
		return
	}

	success := gc.Model.DeleteUserNickname(groupID, request.UserID)
	if !success {
		http.Error(w, "Failed to delete nickname", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Nickname deleted successfully",
		"group_id": groupID,
		"user_id":  request.UserID,
		"nickname": nil,
	})
}

func (gc *GroupController) GetGroupActivities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["id"]
	
	requesterID := r.URL.Query().Get("requester_id")
	if requesterID == "" {
		http.Error(w, "Missing requester_id", http.StatusBadRequest)
		return
	}

	_, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, requesterID) {
		http.Error(w, "Forbidden: Only group members can view group activities", http.StatusForbidden)
		return
	}

	activities, exists := gc.Model.GetGroupActivities(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"group_id":        groupID,
		"activities":      activities,
		"activity_count":  len(activities),
	})
}
