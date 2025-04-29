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
}

func RegisterActivityRoutes(r *mux.Router, activityController *controllers.ActivityController) {
	r.HandleFunc("/activities/{id}", activityController.GetActivity).Methods("GET")
	r.HandleFunc("/activities", activityController.CreateActivity).Methods("POST")
	r.HandleFunc("/activities/{id}", activityController.UpdateActivity).Methods("PUT")
	r.HandleFunc("/activities/{id}", activityController.DeleteActivity).Methods("DELETE")
}
