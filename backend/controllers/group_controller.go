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
	group, exists := gc.Model.GetGroupByID(groupID)
	if !exists {
		http.Error(w, "Group not found", http.StatusNotFound)
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
