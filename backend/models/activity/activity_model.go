package activity

type ActivityModel interface {
	GetActivityByID(id int) (Activity, bool)
	GetActivitiesByCreatorID(creatorID int) []Activity
	CreateActivity(a Activity) Activity
	UpdateActivity(id int, updates map[string]interface{}) (Activity, bool)
	DeleteActivity(id int) bool
}

// DefaultActivityModel must be set in main.go after DB initialization
var DefaultActivityModel ActivityModel
