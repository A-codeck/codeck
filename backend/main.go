package main

import (
	"log"
	"net/http"

	"backend/controllers"
	"backend/models/group"
	"backend/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	groupController := controllers.NewGroupController(group.DefaultGroupModel)

	routes.RegisterRoutes(r, groupController)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
