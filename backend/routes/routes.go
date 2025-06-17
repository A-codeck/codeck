package routes

import (
	"backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterGroupRoutes(r *mux.Router, groupController *controllers.GroupController) {
	r.HandleFunc("/groups/{id}", groupController.GetGroup).Methods("GET")
	r.HandleFunc("/groups", groupController.CreateGroup).Methods("POST")
	r.HandleFunc("/groups/{id}", groupController.UpdateGroup).Methods("PUT")
	r.HandleFunc("/groups/{id}", groupController.DeleteGroup).Methods("DELETE")
	r.HandleFunc("/groups/{id}/members", groupController.GetGroupMembers).Methods("GET")
	r.HandleFunc("/groups/{id}/members", groupController.AddUserToGroup).Methods("POST")
	r.HandleFunc("/groups/{id}/members", groupController.RemoveUserFromGroup).Methods("DELETE")
	r.HandleFunc("/groups/{id}/members/nickname", groupController.SetUserNickname).Methods("PUT")
	r.HandleFunc("/groups/{id}/members/nickname", groupController.DeleteUserNickname).Methods("DELETE")
	r.HandleFunc("/groups/{id}/activities", groupController.GetGroupActivities).Methods("GET")
	r.HandleFunc("/groups/{id}/invites", groupController.CreateInviteLink).Methods("POST")
	r.HandleFunc("/groups/{id}/invites", groupController.GetGroupInvites).Methods("GET")
	r.HandleFunc("/invites/{invite_code}/join", groupController.JoinGroupByInvite).Methods("POST")
	r.HandleFunc("/invites/{invite_code}/deactivate", groupController.DeactivateInvite).Methods("DELETE")
}

func RegisterActivityRoutes(r *mux.Router, activityController *controllers.ActivityController) {
	r.HandleFunc("/activities/{id}", activityController.GetActivity).Methods("GET")
	r.HandleFunc("/activities", activityController.CreateActivity).Methods("POST")
	r.HandleFunc("/activities/{id}", activityController.UpdateActivity).Methods("PUT")
	r.HandleFunc("/activities/{id}", activityController.DeleteActivity).Methods("DELETE")
}

func RegisterUserRoutes(r *mux.Router, userController *controllers.UserController) {
	r.HandleFunc("/users/{id}", userController.GetUser).Methods("GET")
	r.HandleFunc("/users", userController.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}/activities", userController.GetUserActivities).Methods("GET")
}

func RegisterLoginRoutes(r *mux.Router, loginController *controllers.LoginController) {
	r.HandleFunc("/login", loginController.Login).Methods("POST")
}
