package controllers

import (
	"encoding/json"
	"net/http"

	"backend/models/responses"
	"backend/models/user"
)

type LoginController struct {
	Model user.UserModel
}

// swagger imports (used in annotations)
var (
	_ = responses.ErrorResponse{}
)

func NewLoginController(model user.UserModel) *LoginController {
	return &LoginController{Model: model}
}

// Login godoc
// @Summary Authenticate user
// @Description Authenticate user with email and password, returns user data and token
// @Tags authentication
// @Accept json
// @Produce json
// @Param login body responses.LoginRequest true "Login credentials"
// @Success 200 {object} responses.LoginResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 401 {object} responses.ErrorResponse
// @Router /login [post]
func (lc *LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	user, valid := lc.Model.ValidateCredentials(loginRequest.Email, loginRequest.Password)
	if !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"user":  user,
		"token": "dummy-jwt-token-" + user.ID,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
