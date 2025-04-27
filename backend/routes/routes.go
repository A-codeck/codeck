package routes

import (
	"backend/controllers"

	"github.com/gorilla/mux"
)

func RegisterRoutes(r *mux.Router, groupController *controllers.GroupController) {
	r.HandleFunc("/groups/{id}", groupController.GetGroup).Methods("GET")
	r.HandleFunc("/groups", groupController.CreateGroup).Methods("POST")
	r.HandleFunc("/groups/{id}", groupController.UpdateGroup).Methods("PUT")
	r.HandleFunc("/groups/{id}", groupController.DeleteGroup).Methods("DELETE")
}
