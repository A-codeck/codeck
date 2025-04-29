package main

import (
	"log"
	"net/http"

	"backend/controllers"

	"backend/models/activity"
	"backend/models/group"

	"backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	groupController := controllers.NewGroupController(group.DefaultGroupModel)
	activityController := controllers.NewActivityController(activity.DefaultActivityModel)

	routes.RegisterGroupRoutes(r, groupController)
	routes.RegisterActivityRoutes(r, activityController)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
