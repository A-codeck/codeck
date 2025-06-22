package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func (ac *ActivityController) CreateActivity(w http.ResponseWriter, r *http.Request) {
	var raw map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if dateStr, ok := raw["date"].(string); ok && len(dateStr) == 10 {
		raw["date"] = dateStr + "T00:00:00Z"
	}
	fixed, _ := json.Marshal(raw)
	var activity activity.Activity
	if err := json.Unmarshal(fixed, &activity); err != nil {
		log.Printf("Failed to decode request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if activity.Title == "" || activity.Date.IsZero() {
		log.Println("Missing required fields in activity creation")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	createdActivity := ac.Model.CreateActivity(activity)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdActivity)
}

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
