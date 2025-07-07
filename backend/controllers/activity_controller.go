package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"backend/models/activity"
	"backend/models/group"
	"backend/models/responses"
	"backend/utils"

	"github.com/gorilla/mux"
)

type ActivityController struct {
	Model      activity.ActivityModel
	GroupModel group.GroupModel
}

// swagger imports (used in annotations)
var (
	_ = responses.ErrorResponse{}
)

func NewActivityController(model activity.ActivityModel, groupModel group.GroupModel) *ActivityController {
	return &ActivityController{Model: model, GroupModel: groupModel}
}

// GetActivity godoc
// @Summary Get activity by ID
// @Description Get activity information by activity ID
// @Tags activities
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Success 200 {object} activity.Activity
// @Failure 404 {object} responses.ErrorResponse
// @Router /activities/{id} [get]
func (ac *ActivityController) GetActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		log.Printf("Invalid activity id: %v", err)
		http.Error(w, "Invalid activity id", http.StatusBadRequest)
		return
	}
	activity, exists := ac.Model.GetActivityByID(activityID)
	if !exists {
		log.Printf("Activity not found: id=%d", activityID)
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activity)
}

// CreateActivity godoc
// @Summary Create a new activity with image upload and mandatory group selection
// @Description Create a new activity with mandatory image upload, title, date, and at least one group
// @Tags activities
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Activity Title"
// @Param description formData string false "Activity Description"
// @Param date formData string true "Activity Date (YYYY-MM-DD)"
// @Param group_ids formData string true "Comma-separated Group IDs (at least one required)"
// @Param creator_id formData int true "Creator ID"
// @Param image formData file true "Activity Image"
// @Success 201 {object} activity.Activity
// @Failure 400 {object} responses.ErrorResponse
// @Router /activities [post]
func (ac *ActivityController) CreateActivity(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (max 32MB)
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		log.Printf("Failed to parse multipart form: %v", err)
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get form values
	title := r.FormValue("title")
	description := r.FormValue("description")
	date := r.FormValue("date")
	groupIDsStr := r.FormValue("group_ids")
	creatorIDStr := r.FormValue("creator_id")

	// Validate required fields
	if title == "" || date == "" || creatorIDStr == "" || groupIDsStr == "" {
		log.Println("Missing required fields in activity creation")
		http.Error(w, "Missing required fields: title, date, creator_id, and group_ids are required", http.StatusBadRequest)
		return
	}

	// Parse creator_id
	creatorID, err := strconv.Atoi(creatorIDStr)
	if err != nil {
		log.Printf("Invalid creator_id: %v", err)
		http.Error(w, "Invalid creator_id", http.StatusBadRequest)
		return
	}

	// Parse group_ids (comma-separated)
	groupIDStrs := strings.Split(groupIDsStr, ",")
	var groupIDs []int
	for _, idStr := range groupIDStrs {
		if strings.TrimSpace(idStr) == "" {
			continue
		}
		groupID, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			log.Printf("Invalid group_id: %v", err)
			http.Error(w, "Invalid group_id format", http.StatusBadRequest)
			return
		}
		groupIDs = append(groupIDs, groupID)
	}

	// Validate at least one group is provided
	if len(groupIDs) == 0 {
		log.Println("No groups provided for activity")
		http.Error(w, "At least one group is required", http.StatusBadRequest)
		return
	}

	// Get uploaded image file
	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("No image file provided: %v", err)
		http.Error(w, "Image file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the image
	uploadResult, err := utils.SaveImage(file, header)
	if err != nil {
		log.Printf("Failed to save image: %v", err)
		http.Error(w, fmt.Sprintf("Failed to save image: %v", err), http.StatusBadRequest)
		return
	}

	// Fix date format if needed
	if len(date) == 10 {
		date = date + "T00:00:00Z"
	}

	// Create activity struct
	activityData := map[string]interface{}{
		"title":          title,
		"description":    description,
		"activity_image": uploadResult.URL,
		"date":           date,
		"creator_id":     creatorID,
	}

	// Convert to JSON and back to handle date parsing
	jsonData, _ := json.Marshal(activityData)
	var activityObj activity.Activity
	if err := json.Unmarshal(jsonData, &activityObj); err != nil {
		log.Printf("Failed to create activity object: %v", err)
		// Clean up uploaded file on error
		utils.DeleteImage(uploadResult.Filename)
		http.Error(w, "Invalid activity data", http.StatusBadRequest)
		return
	}

	// Validate all groups exist
	for _, groupID := range groupIDs {
		_, groupExists := ac.GroupModel.GetGroupByID(groupID)
		if !groupExists {
			log.Printf("Group not found: id=%d", groupID)
			// Clean up uploaded file on error
			utils.DeleteImage(uploadResult.Filename)
			http.Error(w, fmt.Sprintf("Group not found: %d", groupID), http.StatusNotFound)
			return
		}
	}

	// Create activity with groups in database
	createdActivity := ac.Model.CreateActivityWithGroups(activityObj, groupIDs)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdActivity)
}

// UpdateActivity godoc
// @Summary Update an existing activity
// @Description Update activity information (title cannot be updated)
// @Tags activities
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Param activity body responses.ActivityUpdateRequest true "Activity update data"
// @Success 200 {object} activity.Activity
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /activities/{id} [put]
func (ac *ActivityController) UpdateActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		log.Printf("Invalid activity id: %v", err)
		http.Error(w, "Invalid activity id", http.StatusBadRequest)
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

	updatedActivity, exists := ac.Model.UpdateActivity(activityID, updates)
	if !exists {
		log.Printf("Activity not found: id=%d", activityID)
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedActivity)
}

// DeleteActivity godoc
// @Summary Delete an activity
// @Description Delete an activity (only creator can delete)
// @Tags activities
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Param request body responses.ActivityDeleteRequest true "Delete request with creator_id"
// @Success 204 "No Content"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /activities/{id} [delete]
// DeleteActivity godoc
// @Summary Delete an activity
// @Description Delete an activity (only creator can delete)
// @Tags activities
// @Accept json
// @Produce json
// @Param id path string true "Activity ID"
// @Param request body responses.ActivityDeleteRequest true "Delete request with creator_id"
// @Success 204 "No Content"
// @Failure 400 {object} responses.ErrorResponse
// @Failure 403 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Router /activities/{id} [delete]
func (ac *ActivityController) DeleteActivity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	activityIDStr := vars["id"]
	activityID, err := strconv.Atoi(activityIDStr)
	if err != nil {
		log.Printf("Invalid activity id: %v", err)
		http.Error(w, "Invalid activity id", http.StatusBadRequest)
		return
	}
	activity, exists := ac.Model.GetActivityByID(activityID)
	if !exists {
		log.Printf("Activity not found: id=%d", activityID)
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Failed to decode request payload")
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	creatorIDFloat, ok := request["creator_id"].(float64)
	if !ok {
		log.Println("Missing creator_id in delete activity request")
		http.Error(w, "Missing creator_id", http.StatusBadRequest)
		return
	}
	creatorID := int(creatorIDFloat)

	if creatorID != activity.CreatorID {
		log.Printf("Forbidden: creator_id=%d does not match activity.CreatorID=%d", creatorID, activity.CreatorID)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if !ac.Model.DeleteActivity(activityID) {
		log.Printf("Activity not found for delete: id=%d", activityID)
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserFeed godoc
// @Summary Get user activity feed
// @Description Get activities from all groups the user is a member of
// @Tags activities
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} activity.Activity
// @Failure 400 {object} responses.ErrorResponse
// @Router /activities/feed [get]
func (ac *ActivityController) GetUserFeed(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	log.Printf("Getting feed for user %s", userIDStr)
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Get all groups the user is a member of
	userGroups := ac.GroupModel.GetUserGroups(userID)

	var groupIDs []int
	groupIDs = append(groupIDs, 0) // Include personal activities (group_id = 0)
	for _, group := range userGroups {
		groupIDs = append(groupIDs, group.ID)
	}

	// Get activities from these groups
	activities := ac.Model.GetActivitiesByGroupIDs(groupIDs)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activities)
}

// GetUserFeedWithGroups godoc
// @Summary Get user activity feed with group information
// @Description Get activities from all groups the user is a member of, with group names included
// @Tags activities
// @Accept json
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {array} activity.ActivityWithGroups
// @Failure 400 {object} responses.ErrorResponse
// @Router /activities/feed-with-groups [get]
func (ac *ActivityController) GetUserFeedWithGroups(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "Missing user_id parameter", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user_id: %v", err)
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}

	// Get all groups the user is a member of
	userGroups := ac.GroupModel.GetUserGroups(userID)

	var groupIDs []int
	for _, group := range userGroups {
		groupIDs = append(groupIDs, group.ID)
	}

	// Get activities with group information
	activities := ac.Model.GetActivitiesWithGroupInfo(groupIDs, userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activities)
}
