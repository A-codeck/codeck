package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"backend/models/user"

	"github.com/gorilla/mux"
)

type UserController struct {
	Model user.UserModel
}

func NewUserController(model user.UserModel) *UserController {
	return &UserController{Model: model}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	user, exists := uc.Model.GetUserByID(userID)
	if !exists {
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
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if _, exists := uc.Model.GetUserByEmail(userInput.Email); exists {
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
