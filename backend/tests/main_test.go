package tests

import (
	"os"
	"testing"

	"backend/controllers"
	"backend/models/activity"
	"backend/models/group"
	"backend/routes"

	"github.com/gorilla/mux"
)

var (
	testGroupRouter    *mux.Router
	testGroupModel     *group.InMemoryGroupModel
	testActivityRouter *mux.Router
	testActivityModel  *activity.InMemoryActivityModel
)

func TestMain(m *testing.M) {
	// Initialize the group model and router
	testGroupModel = group.NewInMemoryGroup()
	groupController := controllers.NewGroupController(testGroupModel)
	testGroupRouter = mux.NewRouter()
	routes.RegisterGroupRoutes(testGroupRouter, groupController)

	// Initialize the activity model and router
	testActivityModel = activity.NewInMemoryActivity()
	activityController := controllers.NewActivityController(testActivityModel)
	testActivityRouter = mux.NewRouter()
	routes.RegisterActivityRoutes(testActivityRouter, activityController)

	// Run all tests
	os.Exit(m.Run())
}
