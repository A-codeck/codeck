package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"backend/models/activity"

	"github.com/gorilla/mux"
)

type ActivityController struct {
	Model activity.ActivityModel
}

func NewActivityController(model activity.ActivityModel) *ActivityController {
	return &ActivityController{Model: model}
}

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
