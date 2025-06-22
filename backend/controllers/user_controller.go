package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"backend/models/activity"
	"backend/models/user"

	"github.com/gorilla/mux"
)

type UserController struct {
	Model         user.UserModel
	ActivityModel activity.ActivityModel
}

func NewUserController(model user.UserModel, activityModel activity.ActivityModel) *UserController {
	return &UserController{Model: model, ActivityModel: activityModel}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user id: %v", err)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	user, exists := uc.Model.GetUserByID(userID)
	if !exists {
		log.Printf("User not found: id=%d", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	body, _ := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		http.Error(w, "Invalid request payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if userInput.Email == "" || userInput.Name == "" || userInput.Password == "" {
		log.Println("Missing required fields in user creation")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if _, exists := uc.Model.GetUserByEmail(userInput.Email); exists {
		log.Printf("Email already in use: %s", userInput.Email)
		http.Error(w, "Email already in use", http.StatusConflict)
		return
	}

	user := user.User{
		Email:    userInput.Email,
		Name:     userInput.Name,
		Password: userInput.Password,
	}

	createdUser := uc.Model.CreateUser(user)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdUser)
}

func (uc *UserController) GetUserActivities(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["id"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user id: %v", err)
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	// Check if user exists
	_, exists := uc.Model.GetUserByID(userID)
	if !exists {
		log.Printf("User not found: id=%d", userID)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Get all activities for this user
	activities := uc.ActivityModel.GetActivitiesByCreatorID(userID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activities)
}
