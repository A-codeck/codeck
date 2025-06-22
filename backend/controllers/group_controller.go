package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		log.Printf("Invalid group id: %v", err)
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	requesterIDStr := r.URL.Query().Get("requester_id")
	requesterID, err := strconv.Atoi(requesterIDStr)
	if err != nil {
		log.Printf("Invalid requester_id: %v", err)
		http.Error(w, "Invalid requester_id", http.StatusBadRequest)
		return
	}

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		log.Printf("Group not found: id=%d", groupID)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, requesterID) {
		log.Printf("Forbidden: requester_id=%d is not a member of group_id=%d", requesterID, groupID)
		http.Error(w, "Forbidden: Only group members can view group details", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}

func (gc *GroupController) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if dateStr, ok := raw["end_date"].(string); ok && len(dateStr) == 10 {
		raw["end_date"] = dateStr + "T00:00:00Z"
	}
	fixed, _ := json.Marshal(raw)
	var group group.Group
	if err := json.Unmarshal(fixed, &group); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if group.Name == "" || group.EndDate.IsZero() {
		log.Println("Missing required fields in group creation")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	createdGroup := gc.Model.CreateGroup(group)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdGroup)
}

func (gc *GroupController) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		log.Printf("Invalid group id: %v", err)
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, ok := updates["name"].(string); ok {
		log.Println("Name field cannot be updated")
		http.Error(w, "Name field cannot be updated", http.StatusBadRequest)
		return
	}

	updatedGroup, exists := gc.Model.UpdateGroup(groupID, updates)
	if !exists {
		log.Printf("Group not found: id=%d", groupID)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedGroup)
}

func (gc *GroupController) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		log.Printf("Invalid group id: %v", err)
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}
	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		log.Printf("Group not found: id=%d", groupID)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Failed to decode request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	creatorID, ok := request["creator_id"].(float64)
	if !ok {
		log.Println("Missing creator_id in delete group request")
		http.Error(w, "Missing creator_id", http.StatusBadRequest)
		return
	}

	if int(creatorID) != group.CreatorID {
		log.Printf("Forbidden: creator_id=%d does not match group.CreatorID=%d", int(creatorID), group.CreatorID)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if !gc.Model.DeleteGroup(groupID) {
		log.Printf("Group not found for delete: id=%d", groupID)
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (gc *GroupController) AddUserToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	_, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID int `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 {
		log.Println("Missing user_id in add user to group request")
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	if gc.Model.IsUserInGroup(groupID, request.UserID) {
		log.Printf("User is already a member of group: group_id=%d, user_id=%d", groupID, request.UserID)
		http.Error(w, "User is already a member of this group", http.StatusConflict)
		return
	}

	success := gc.Model.AddUserToGroup(groupID, request.UserID)
	if !success {
		log.Printf("Failed to add user to group: group_id=%d, user_id=%d", groupID, request.UserID)
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
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID      int `json:"user_id"`
		RequesterID int `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 || request.RequesterID == 0 {
		log.Println("Missing user_id or requester_id in remove user from group request")
		http.Error(w, "Missing user_id or requester_id", http.StatusBadRequest)
		return
	}

	if request.RequesterID != group.CreatorID && request.RequesterID != request.UserID {
		log.Printf("Forbidden: requester_id=%d is not allowed to remove user_id=%d from group_id=%d", request.RequesterID, request.UserID, groupID)
		http.Error(w, "Forbidden: Only group creator or the user themselves can remove membership", http.StatusForbidden)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, request.UserID) {
		log.Printf("User is not a member of group: group_id=%d, user_id=%d", groupID, request.UserID)
		http.Error(w, "User is not a member of this group", http.StatusNotFound)
		return
	}

	success := gc.Model.RemoveUserFromGroup(groupID, request.UserID)
	if !success {
		log.Printf("Failed to remove user from group: group_id=%d, user_id=%d", groupID, request.UserID)
		http.Error(w, "Failed to remove user from group", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "User removed from group successfully",
		"group_id": groupID,
		"user_id":  request.UserID,
	})
}

func (gc *GroupController) GetGroupMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	requesterIDStr := r.URL.Query().Get("requester_id")
	requesterID, err := strconv.Atoi(requesterIDStr)
	if err != nil {
		http.Error(w, "Invalid requester_id", http.StatusBadRequest)
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
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		CreatorID int     `json:"creator_id"`
		ExpiresAt *string `json:"expires_at,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.CreatorID == 0 {
		log.Println("Missing creator_id in create invite link request")
		http.Error(w, "Missing creator_id", http.StatusBadRequest)
		return
	}

	if request.CreatorID != group.CreatorID {
		log.Printf("Forbidden: creator_id=%d is not group creator (group.CreatorID=%d)", request.CreatorID, group.CreatorID)
		http.Error(w, "Forbidden: Only group creator can create invite links", http.StatusForbidden)
		return
	}

	invite, success := gc.Model.CreateInviteLink(groupID, request.CreatorID, request.ExpiresAt)
	if !success {
		log.Printf("Failed to create invite link for group_id=%d by creator_id=%d", groupID, request.CreatorID)
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
		UserID      int `json:"user_id"`
		RequesterID int `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 || request.RequesterID == 0 {
		log.Println("Missing user_id or requester_id in join group by invite request")
		http.Error(w, "Missing user_id or requester_id", http.StatusBadRequest)
		return
	}

	invite, exists := gc.Model.GetInviteByCode(inviteCode)
	if !exists || !invite.IsActive {
		log.Printf("Invalid or expired invite code: %s", inviteCode)
		http.Error(w, "Invalid or expired invite code", http.StatusNotFound)
		return
	}

	group, exists := gc.Model.GetGroupByID(invite.GroupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if gc.Model.IsUserInGroup(invite.GroupID, request.UserID) {
		log.Printf("User is already a member of group: group_id=%d, user_id=%d", invite.GroupID, request.UserID)
		http.Error(w, "User is already a member of this group", http.StatusConflict)
		return
	}

	if request.RequesterID != group.CreatorID && request.RequesterID != invite.CreatedBy {
		log.Printf("Forbidden: requester_id=%d is not allowed to join group_id=%d by invite_code=%s", request.RequesterID, invite.GroupID, inviteCode)
		http.Error(w, "Forbidden: Only group creator or invite creator can join group by invite", http.StatusForbidden)
		return
	}

	success := gc.Model.DeactivateInvite(inviteCode)
	if !success {
		log.Printf("Failed to deactivate invite: invite_code=%s", inviteCode)
		http.Error(w, "Failed to deactivate invite", http.StatusInternalServerError)
		return
	}

	success = gc.Model.AddUserToGroup(invite.GroupID, request.UserID)
	if !success {
		log.Printf("Failed to join group: group_id=%d, user_id=%d", invite.GroupID, request.UserID)
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
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

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
		RequesterID int `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.RequesterID == 0 {
		log.Println("Missing requester_id in deactivate invite request")
		http.Error(w, "Missing requester_id", http.StatusBadRequest)
		return
	}

	invite, exists := gc.Model.GetInviteByCode(inviteCode)
	if !exists {
		log.Printf("Invite not found: invite_code=%s", inviteCode)
		http.Error(w, "Invite not found", http.StatusNotFound)
		return
	}

	group, exists := gc.Model.GetGroupByID(invite.GroupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if request.RequesterID != group.CreatorID && request.RequesterID != invite.CreatedBy {
		log.Printf("Forbidden: requester_id=%d is not allowed to deactivate invite_code=%s", request.RequesterID, inviteCode)
		http.Error(w, "Forbidden: Only group creator or invite creator can deactivate invite", http.StatusForbidden)
		return
	}

	success := gc.Model.DeactivateInvite(inviteCode)
	if !success {
		log.Printf("Failed to deactivate invite: invite_code=%s", inviteCode)
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
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID      int     `json:"user_id"`
		RequesterID int     `json:"requester_id"`
		Nickname    *string `json:"nickname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 || request.RequesterID == 0 {
		log.Println("Missing user_id or requester_id in set nickname request")
		http.Error(w, "Missing user_id or requester_id", http.StatusBadRequest)
		return
	}

	if request.RequesterID != request.UserID && request.RequesterID != group.CreatorID {
		log.Printf("Forbidden: requester_id=%d is not allowed to set nickname for user_id=%d in group_id=%d", request.RequesterID, request.UserID, groupID)
		http.Error(w, "Forbidden: Only the user themselves or group creator can set nickname", http.StatusForbidden)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, request.UserID) {
		log.Printf("User is not a member of group: group_id=%d, user_id=%d", groupID, request.UserID)
		http.Error(w, "User is not a member of this group", http.StatusNotFound)
		return
	}

	if request.Nickname != nil && len(*request.Nickname) > 50 {
		log.Printf("Nickname too long: user_id=%d, group_id=%d", request.UserID, groupID)
		http.Error(w, "Nickname cannot be longer than 50 characters", http.StatusBadRequest)
		return
	}

	success := gc.Model.SetUserNickname(groupID, request.UserID, request.Nickname)
	if !success {
		log.Printf("Failed to set nickname: group_id=%d, user_id=%d", groupID, request.UserID)
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
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	var request struct {
		UserID      int `json:"user_id"`
		RequesterID int `json:"requester_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 || request.RequesterID == 0 {
		log.Println("Missing user_id or requester_id in delete nickname request")
		http.Error(w, "Missing user_id or requester_id", http.StatusBadRequest)
		return
	}

	if request.RequesterID != request.UserID && request.RequesterID != group.CreatorID {
		log.Printf("Forbidden: requester_id=%d is not allowed to delete nickname for user_id=%d in group_id=%d", request.RequesterID, request.UserID, groupID)
		http.Error(w, "Forbidden: Only the user themselves or group creator can delete nickname", http.StatusForbidden)
		return
	}

	if !gc.Model.IsUserInGroup(groupID, request.UserID) {
		log.Printf("User is not a member of group: group_id=%d, user_id=%d", groupID, request.UserID)
		http.Error(w, "User is not a member of this group", http.StatusNotFound)
		return
	}

	success := gc.Model.DeleteUserNickname(groupID, request.UserID)
	if !success {
		log.Printf("Failed to delete nickname: group_id=%d, user_id=%d", groupID, request.UserID)
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
	groupIDStr := vars["id"]
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "Invalid group id", http.StatusBadRequest)
		return
	}

	requesterIDStr := r.URL.Query().Get("requester_id")
	requesterID, err := strconv.Atoi(requesterIDStr)
	if err != nil {
		http.Error(w, "Invalid requester_id", http.StatusBadRequest)
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
		"group_id":       groupID,
		"activities":     activities,
		"activity_count": len(activities),
	})
}
