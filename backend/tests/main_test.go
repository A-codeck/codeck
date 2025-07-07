package tests

import (
	"fmt"
	"log"
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
	testDB             *gorm.DB
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

// getTestDBConfig returns the database configuration for tests
func getTestDBConfig() string {
	host := getEnvOrDefault("TEST_DB_HOST", "localhost")
	user := getEnvOrDefault("TEST_DB_USER", "my_usr")
	password := getEnvOrDefault("TEST_DB_PASSWORD", "my_pwd")
	dbname := getEnvOrDefault("TEST_DB_NAME", "codeck_test")
	port := getEnvOrDefault("TEST_DB_PORT", "5433")
	sslmode := getEnvOrDefault("TEST_DB_SSLMODE", "disable")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		host, user, password, dbname, port, sslmode)
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setupTestDB creates and migrates the test database
func setupTestDB() error {
	dsn := getTestDBConfig()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %w\nMake sure the test database is running: docker-compose -f docker-compose.test.yml up -d", err)
	}

	testDB = db

	// Run migrations for test database
	err = db.AutoMigrate(
		&group.Group{},
		&group.GroupMember{},
		&group.GroupInvite{},
		&activity.Activity{},
		&comment.Comment{},
		&user.User{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate test database: %w", err)
	}

	return nil
}

// cleanTestDB removes all data from test database
func cleanTestDB() error {
	if testDB == nil {
		return nil
	}

	// Delete all records in reverse order to respect foreign key constraints
	tables := []interface{}{
		&comment.Comment{},
		&activity.Activity{},
		&group.GroupInvite{},
		&group.GroupMember{},
		&group.Group{},
		&user.User{},
	}

	for _, table := range tables {
		if err := testDB.Unscoped().Where("1 = 1").Delete(table).Error; err != nil {
			return fmt.Errorf("failed to clean table %T: %w", table, err)
		}
	}

	// Reset auto-increment sequences for PostgreSQL
	resetQueries := []string{
		"ALTER SEQUENCE users_id_seq RESTART WITH 1",
		"ALTER SEQUENCE groups_id_seq RESTART WITH 1",
		"ALTER SEQUENCE activities_id_seq RESTART WITH 1",
		"ALTER SEQUENCE comments_id_seq RESTART WITH 1",
	}

	for _, query := range resetQueries {
		if err := testDB.Exec(query).Error; err != nil {
			// Log the error but don't fail - sequence might not exist yet
			fmt.Printf("Warning: failed to reset sequence with query '%s': %v\n", query, err)
		}
	}

	return nil
}

// closeTestDB closes the test database connection
func closeTestDB() error {
	if testDB == nil {
		return nil
	}

	sqlDB, err := testDB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

func TestMain(m *testing.M) {
	// Setup test database
	if err := setupTestDB(); err != nil {
		log.Fatalf("Failed to setup test database: %v", err)
	}

	// Clean database before tests
	if err := cleanTestDB(); err != nil {
		log.Fatalf("Failed to clean test database: %v", err)
	}

	// Initialize models with test database
	testGroupModel = group.NewGormGroupModel(testDB)
	testActivityModel = activity.NewGormActivityModel(testDB)
	testUserModel = user.NewGormUserModel(testDB)
	testCommentModel = comment.NewGormCommentModel(testDB)

	// Initialize controllers and routers
	groupController := controllers.NewGroupController(testGroupModel, testActivityModel, testUserModel)
	testGroupRouter = mux.NewRouter()
	routes.RegisterGroupRoutes(testGroupRouter, groupController)

	activityController := controllers.NewActivityController(testActivityModel, testGroupModel)
	testActivityRouter = mux.NewRouter()
	routes.RegisterActivityRoutes(testActivityRouter, activityController)

	userController := controllers.NewUserController(testUserModel, testActivityModel, testGroupModel)
	testUserRouter = mux.NewRouter()
	routes.RegisterUserRoutes(testUserRouter, userController)

	loginController := controllers.NewLoginController(testUserModel)
	testLoginRouter = mux.NewRouter()
	routes.RegisterLoginRoutes(testLoginRouter, loginController)

	commentController := controllers.NewCommentController(testCommentModel, testActivityModel, testGroupModel)
	testCommentRouter = mux.NewRouter()
	routes.RegisterCommentRoutes(testCommentRouter, commentController)

	// Run tests
	code := m.Run()

	// Cleanup after tests
	if err := cleanTestDB(); err != nil {
		log.Printf("Warning: Failed to clean test database after tests: %v", err)
	}

	if err := closeTestDB(); err != nil {
		log.Printf("Warning: Failed to close test database connection: %v", err)
	}

	os.Exit(code)
}

// CleanupBetweenTests can be called at the beginning of each test to ensure a clean state
func CleanupBetweenTests(t *testing.T) {
	t.Helper()
	if err := cleanTestDB(); err != nil {
		t.Fatalf("Failed to cleanup database between tests: %v", err)
	}
}

// SeedAllTestData seeds all necessary test data in the correct order
func SeedAllTestData(t *testing.T) {
	t.Helper()

	// Clean first
	CleanupBetweenTests(t)

	// Seed in dependency order: users first, then groups, then activities, then comments
	testUserModel.SeedDefaultData()
	testGroupModel.SeedDefaultData()
	testActivityModel.SeedDefaultData()
	testCommentModel.SeedDefaultData()
}
