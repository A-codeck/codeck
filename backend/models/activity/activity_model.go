package activity

type ActivityModel interface {
	GetActivityByID(id string) (Activity, bool)
	GetActivitiesByCreatorID(creatorID string) []Activity
	GetActivitiesByGroupIDs(groupIDs []string) []Activity
	CreateActivity(group Activity) Activity
	UpdateActivity(id string, updates map[string]interface{}) (Activity, bool)
	DeleteActivity(id string) bool
}

var DefaultActivityModel ActivityModel = NewInMemoryActivity()
