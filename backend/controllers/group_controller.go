package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"backend/models/activity"
	"backend/models/group"
	"backend/models/responses"

	"github.com/gorilla/mux"
)

type GroupController struct {
	Model         group.GroupModel
	ActivityModel activity.ActivityModel
}

// swagger imports (used in annotations)
var (
	_ = activity.Activity{}
	_ = responses.ErrorResponse{}
)

func NewGroupController(model group.GroupModel, activityModel activity.ActivityModel) *GroupController {
	return &GroupController{Model: model, ActivityModel: activityModel}
}

// GetGroup godoc
// @Summary Get group details
// @Description Get group information (members only)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param requester_id query int true "Requester User ID"
// @Success 200 {object} group.Group
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id} [get]
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

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new group with name, end date, and optional image/description
// @Tags groups
// @Accept json
// @Produce json
// @Param group body responses.GroupCreateRequest true "Group creation data"
// @Success 201 {object} group.Group
// @Failure 400 {object} responses.ErrorResponse
// @Router /groups [post]
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

// UpdateGroup godoc
// @Summary Update an existing group
// @Description Update group information (name cannot be updated)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param group body responses.GroupUpdateRequest true "Group update data"
// @Success 200 {object} group.Group
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id} [put]
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

// DeleteGroup godoc
// @Summary Delete a group
// @Description Delete a group (only creator can delete)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body responses.GroupDeleteRequest true "Delete request with creator_id"
// @Success 204 "No Content"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id} [delete]
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

// AddUserToGroup godoc
// @Summary Add user to group
// @Description Add a user to a group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body responses.AddUserToGroupRequest true "Add user request"
// @Success 201 {object} responses.AddUserToGroupResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Router /groups/{id}/members [post]
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

	// Check if user exists (you'll need to inject UserModel into GroupController)
	// For now, we'll skip this check since we don't have access to UserModel here

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

// RemoveUserFromGroup godoc
// @Summary Remove user from group
// @Description Remove a user from a group (only group creator or the user themselves can remove)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body responses.RemoveUserFromGroupRequest true "Remove user request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id}/members [delete]
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

// GetGroupMembers godoc
// @Summary Get group members
// @Description Get all members of a group (members only)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param requester_id query int true "Requester User ID"
// @Success 200 {object} responses.GroupMembersResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id}/members [get]
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

// CreateInviteLink godoc
// @Summary Create group invite link
// @Description Create an invite link for the group (only group creator can create invites)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body responses.CreateInviteRequest true "Create invite request"
// @Success 201 {object} group.GroupInvite
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id}/invites [post]
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

// JoinGroupByInvite godoc
// @Summary Join group by invite code
// @Description Join a group using an invite code
// @Tags groups
// @Accept json
// @Produce json
// @Param invite_code path string true "Invite Code"
// @Param request body responses.JoinGroupRequest true "Join group request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Router /invites/{invite_code}/join [post]
func (gc *GroupController) JoinGroupByInvite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	inviteCode := vars["invite_code"]

	var request struct {
		UserID   int     `json:"user_id"`
		Nickname *string `json:"nickname,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if request.UserID == 0 {
		log.Println("Missing user_id in join group by invite request")
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	invite, exists := gc.Model.GetInviteByCode(inviteCode)
	if !exists || !invite.IsActive {
		log.Printf("Invalid or expired invite code: %s", inviteCode)
		http.Error(w, "Invalid or expired invite code", http.StatusNotFound)
		return
	}

	_, exists = gc.Model.GetGroupByID(invite.GroupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
		return
	}

	if gc.Model.IsUserInGroup(invite.GroupID, request.UserID) {
		log.Printf("User is already a member of group: group_id=%d, user_id=%d", invite.GroupID, request.UserID)
		http.Error(w, "User is already a member of this group", http.StatusConflict)
		return
	}

	success := gc.Model.AddUserToGroup(invite.GroupID, request.UserID)
	if !success {
		log.Printf("Failed to join group: group_id=%d, user_id=%d", invite.GroupID, request.UserID)
		http.Error(w, "Failed to join group", http.StatusInternalServerError)
		return
	}

	// If nickname is provided, set it for the user
	if request.Nickname != nil {
		if len(*request.Nickname) > 50 {
			log.Printf("Nickname too long: user_id=%d, group_id=%d", request.UserID, invite.GroupID)
			http.Error(w, "Nickname cannot be longer than 50 characters", http.StatusBadRequest)
			return
		}
		nickSuccess := gc.Model.SetUserNickname(invite.GroupID, request.UserID, request.Nickname)
		if !nickSuccess {
			log.Printf("Failed to set nickname: group_id=%d, user_id=%d", invite.GroupID, request.UserID)
			http.Error(w, "Failed to set nickname", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":     "Successfully joined group",
		"group_id":    invite.GroupID,
		"user_id":     request.UserID,
		"nickname":    request.Nickname,
		"invite_code": inviteCode,
	}

	json.NewEncoder(w).Encode(response)
}

// GetGroupInvites godoc
// @Summary Get group invites
// @Description Get all invites for a group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Success 200 {array} group.GroupInvite
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id}/invites [get]
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

// DeactivateInvite godoc
// @Summary Deactivate group invite
// @Description Deactivate an invite link (only group creator can deactivate)
// @Tags groups
// @Accept json
// @Produce json
// @Param invite_code path string true "Invite Code"
// @Param request body responses.DeactivateInviteRequest true "Deactivate invite request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /invites/{invite_code}/deactivate [delete]
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

// SetUserNickname godoc
// @Summary Set user nickname in group
// @Description Set or update a user's nickname in a group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body responses.SetNicknameRequest true "Set nickname request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id}/members/nickname [put]
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

// DeleteUserNickname godoc
// @Summary Delete user nickname in group
// @Description Remove a user's nickname in a group
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "Group ID"
// @Param request body responses.DeleteNicknameRequest true "Delete nickname request"
// @Success 200 {object} responses.SuccessResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id}/members/nickname [delete]
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

// GetGroupActivities godoc
// @Summary Get group activities
// @Description Get all activities for a group (members only)
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param requester_id query int true "Requester User ID"
// @Success 200 {object} responses.GroupActivitiesResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /groups/{id}/activities [get]
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

	// Get activities for this group
	activities := gc.ActivityModel.GetActivitiesByGroupID(groupID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"group_id":       groupID,
		"activities":     activities,
		"activity_count": len(activities),
	})
}
