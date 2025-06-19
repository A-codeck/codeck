package tests

import (
	"os"
	"testing"

	"backend/controllers"
	"backend/models/activity"
	"backend/models/comment"
	"backend/models/group"
	"backend/models/user"
	"backend/routes"

	"github.com/gorilla/mux"
)

var (
	testGroupRouter *mux.Router
	testGroupModel  *group.InMemoryGroupModel

	testActivityRouter *mux.Router
	testActivityModel  *activity.InMemoryActivityModel

	testUserRouter *mux.Router
	testUserModel  *user.InMemoryUserModel

	testCommentRouter *mux.Router
	testCommentModel  *comment.InMemoryCommentModel

	testLoginRouter *mux.Router
)

func TestMain(m *testing.M) {
	testGroupModel = group.NewInMemoryGroup()
	groupController := controllers.NewGroupController(testGroupModel)
	testGroupRouter = mux.NewRouter()
	routes.RegisterGroupRoutes(testGroupRouter, groupController)

	testActivityModel = activity.NewInMemoryActivity()
	activityController := controllers.NewActivityController(testActivityModel)
	testActivityRouter = mux.NewRouter()
	routes.RegisterActivityRoutes(testActivityRouter, activityController)

	testUserModel = user.NewInMemoryUser()
	userController := controllers.NewUserController(testUserModel, testActivityModel)
	testUserRouter = mux.NewRouter()
	routes.RegisterUserRoutes(testUserRouter, userController)

	loginController := controllers.NewLoginController(testUserModel)
	testLoginRouter = mux.NewRouter()
	routes.RegisterLoginRoutes(testLoginRouter, loginController)

	testCommentModel = comment.NewInMemoryComment()
	commentController := controllers.NewCommentController(testCommentModel, testActivityModel, testGroupModel)
	testCommentRouter = mux.NewRouter()
	routes.RegisterCommentRoutes(testCommentRouter, commentController)

	os.Exit(m.Run())
}
