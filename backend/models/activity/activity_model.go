package activity

type ActivityModel interface {
	GetActivityByID(id int) (Activity, bool)
	GetActivitiesByCreatorID(creatorID int) []Activity
	GetActivitiesByGroupID(groupID int) []Activity
	GetActivitiesByGroupIDs(groupIDs []int) []Activity
	CreateActivity(a Activity) Activity
	CreateActivityWithGroups(a Activity, groupIDs []int) Activity
	UpdateActivity(id int, updates map[string]interface{}) (Activity, bool)
	DeleteActivity(id int) bool
	GetActivityGroups(activityID int) []int
	SetActivityGroups(activityID int, groupIDs []int) error
	GetActivitiesWithGroupInfo(groupIDs []int, userID int) []ActivityWithGroups
}

// DefaultActivityModel must be set in main.go after DB initialization
var DefaultActivityModel ActivityModel
