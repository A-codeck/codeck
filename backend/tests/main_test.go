package tests

import (
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/controllers"
	"backend/models/activity"
	"backend/models/comment"
	"backend/models/group"
	"backend/models/user"
	"backend/routes"

	"github.com/gorilla/mux"
)

var (
	testGroupRouter    *mux.Router
	testGroupModel     *group.GormGroupModel
	testActivityRouter *mux.Router
	testActivityModel  *activity.GormActivityModel
	testUserRouter     *mux.Router
	testUserModel      *user.GormUserModel
	testCommentRouter  *mux.Router
	testCommentModel   *comment.GormCommentModel
	testLoginRouter    *mux.Router
)

func TestMain(m *testing.M) {
	dsn := "host=0.0.0.0 user=my_usr password=my_pwd dbname=codeck port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&group.Group{}, &activity.Activity{}, &comment.Comment{}, &user.User{})

	testGroupModel = group.NewGormGroupModel(db)
	groupController := controllers.NewGroupController(testGroupModel)
	testGroupRouter = mux.NewRouter()
	routes.RegisterGroupRoutes(testGroupRouter, groupController)

	testActivityModel = activity.NewGormActivityModel(db)
	activityController := controllers.NewActivityController(testActivityModel)
	testActivityRouter = mux.NewRouter()
	routes.RegisterActivityRoutes(testActivityRouter, activityController)

	testUserModel = user.NewGormUserModel(db)
	userController := controllers.NewUserController(testUserModel, testActivityModel)
	testUserRouter = mux.NewRouter()
	routes.RegisterUserRoutes(testUserRouter, userController)

	loginController := controllers.NewLoginController(testUserModel)
	testLoginRouter = mux.NewRouter()
	routes.RegisterLoginRoutes(testLoginRouter, loginController)

	testCommentModel = comment.NewGormCommentModel(db)
	commentController := controllers.NewCommentController(testCommentModel, testActivityModel, testGroupModel)
	testCommentRouter = mux.NewRouter()
	routes.RegisterCommentRoutes(testCommentRouter, commentController)

	os.Exit(m.Run())
}
