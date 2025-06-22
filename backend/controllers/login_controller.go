package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"backend/models/user"
)

type LoginController struct {
	Model user.UserModel
}

func NewLoginController(model user.UserModel) *LoginController {
	return &LoginController{Model: model}
}

func (lc *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Printf("Failed to decode login request payload: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		log.Println("Missing email or password in login request")
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, valid := lc.Model.ValidateCredentials(loginRequest.Email, loginRequest.Password)
	if !valid {
		log.Printf("Invalid credentials for email: %s", loginRequest.Email)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": "dummy-jwt-token-" + strconv.Itoa(user.ID),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
