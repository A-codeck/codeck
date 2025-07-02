package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/models/activity"
	"backend/models/responses"

	"github.com/gorilla/mux"
)

type ActivityController struct {
	Model activity.ActivityModel
}

// swagger imports (used in annotations)
var (
	_ = responses.ErrorResponse{}
)

func NewActivityController(model activity.ActivityModel) *ActivityController {
	return &ActivityController{Model: model}
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
	activityID := vars["id"]
	activity, exists := ac.Model.GetActivityByID(activityID)
	if !exists {
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activity)
}

// CreateActivity godoc
// @Summary Create a new activity
// @Description Create a new activity with title, date, and optional image/description
// @Tags activities
// @Accept json
// @Produce json
// @Param activity body responses.ActivityCreateRequest true "Activity creation data"
// @Success 201 {object} activity.Activity
// @Failure 400 {object} responses.ErrorResponse
// @Router /activities [post]
func (ac *ActivityController) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var activity activity.Activity
	if err := json.NewDecoder(r.Body).Decode(&activity); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if activity.Title == "" || activity.Date == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	createdActivity := ac.Model.CreateActivity(activity)

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
	activityID := vars["id"]

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if _, ok := updates["name"].(string); ok {
		http.Error(w, "Name field cannot be updated", http.StatusBadRequest)
		return
	}

	updatedActivity, exists := ac.Model.UpdateActivity(activityID, updates)
	if !exists {
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
	activityID := vars["id"]
	activity, exists := ac.Model.GetActivityByID(activityID)
	if !exists {
		http.Error(w, "Activity not found", http.StatusNotFound)
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

	if creatorID != activity.CreatorID {
		// Activity exists but creator_id is invalid
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if !ac.Model.DeleteActivity(activityID) {
		http.Error(w, "Activity not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
