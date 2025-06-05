package main

import (
	"log"
	"net/http"

	"backend/controllers"

	"backend/models/activity"
	"backend/models/group"

	"backend/routes"

	"github.com/gorilla/mux"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := mux.NewRouter()

	groupController := controllers.NewGroupController(group.DefaultGroupModel)
	activityController := controllers.NewActivityController(activity.DefaultActivityModel)

	routes.RegisterGroupRoutes(r, groupController)
	routes.RegisterActivityRoutes(r, activityController)

	log.Println("Joao is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	log.Println("Trying to migrate")
	dsn := "host=db user=my_usr password=my_pwd dbname=codeck port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	if err := db.AutoMigrate(&group.Group{}, &activity.Activity{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
