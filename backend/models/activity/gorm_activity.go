package activity

import (
	"time"

	"gorm.io/gorm"
)

type GormActivityModel struct {
	db *gorm.DB
}

func NewGormActivityModel(db *gorm.DB) *GormActivityModel {
	return &GormActivityModel{db: db}
}

func (m *GormActivityModel) GetActivityByID(id int) (Activity, bool) {
	var a Activity
	if err := m.db.First(&a, "id = ?", id).Error; err != nil {
		return Activity{}, false
	}
	// Load associated groups
	a.Groups = m.GetActivityGroups(a.ID)
	return a, true
}

func (m *GormActivityModel) GetActivitiesByCreatorID(creatorID int) []Activity {
	var list []Activity
	m.db.Where("creator_id = ?", creatorID).Find(&list)
	// Load groups for each activity
	for i := range list {
		list[i].Groups = m.GetActivityGroups(list[i].ID)
	}
	return list
}

func (m *GormActivityModel) GetActivitiesByGroupID(groupID int) []Activity {
	var activityIDs []int
	m.db.Model(&ActivityGroup{}).
		Where("group_id = ?", groupID).
		Pluck("activity_id", &activityIDs)

	if len(activityIDs) == 0 {
		return []Activity{}
	}

	var list []Activity
	m.db.Where("id IN ?", activityIDs).Find(&list)
	// Load groups for each activity
	for i := range list {
		list[i].Groups = m.GetActivityGroups(list[i].ID)
	}
	return list
}

func (m *GormActivityModel) GetActivitiesByGroupIDs(groupIDs []int) []Activity {
	var activityIDs []int
	m.db.Model(&ActivityGroup{}).
		Where("group_id IN ?", groupIDs).
		Pluck("activity_id", &activityIDs)

	if len(activityIDs) == 0 {
		return []Activity{}
	}

	var list []Activity
	m.db.Where("id IN ?", activityIDs).Find(&list)
	// Load groups for each activity
	for i := range list {
		list[i].Groups = m.GetActivityGroups(list[i].ID)
	}
	return list
}

func (m *GormActivityModel) CreateActivity(a Activity) Activity {
	m.db.Create(&a)
	return a
}

func (m *GormActivityModel) CreateActivityWithGroups(a Activity, groupIDs []int) Activity {
	// Create activity first
	m.db.Create(&a)

	// Then create group associations
	for _, groupID := range groupIDs {
		activityGroup := ActivityGroup{
			ActivityID: a.ID,
			GroupID:    groupID,
		}
		m.db.Create(&activityGroup)
	}

	// Load groups back to the activity
	a.Groups = groupIDs
	return a
}

func (m *GormActivityModel) UpdateActivity(id int, updates map[string]interface{}) (Activity, bool) {
	var a Activity
	if err := m.db.First(&a, "id = ?", id).Error; err != nil {
		return Activity{}, false
	}
	m.db.Model(&a).Updates(updates)
	// Load associated groups
	a.Groups = m.GetActivityGroups(a.ID)
	return a, true
}

func (m *GormActivityModel) DeleteActivity(id int) bool {
	// Delete group associations first
	m.db.Where("activity_id = ?", id).Delete(&ActivityGroup{})
	// Then delete the activity
	if err := m.db.Delete(&Activity{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

func (m *GormActivityModel) GetActivityGroups(activityID int) []int {
	var groupIDs []int
	m.db.Model(&ActivityGroup{}).
		Where("activity_id = ?", activityID).
		Pluck("group_id", &groupIDs)
	return groupIDs
}

func (m *GormActivityModel) SetActivityGroups(activityID int, groupIDs []int) error {
	// First delete existing associations
	m.db.Where("activity_id = ?", activityID).Delete(&ActivityGroup{})

	// Then create new associations
	for _, groupID := range groupIDs {
		activityGroup := ActivityGroup{
			ActivityID: activityID,
			GroupID:    groupID,
		}
		if err := m.db.Create(&activityGroup).Error; err != nil {
			return err
		}
	}
	return nil
}

func (m *GormActivityModel) GetActivitiesWithGroupInfo(groupIDs []int, userID int) []ActivityWithGroups {
	// Get distinct activities from the user's groups
	var activityIDs []int
	m.db.Model(&ActivityGroup{}).
		Where("group_id IN ?", groupIDs).
		Group("activity_id").
		Pluck("activity_id", &activityIDs)

	if len(activityIDs) == 0 {
		return []ActivityWithGroups{}
	}

	var activities []Activity
	m.db.Where("id IN ?", activityIDs).Find(&activities)

	var result []ActivityWithGroups
	for _, activity := range activities {
		// Get group names for this activity where the user is a member
		var groupNames []string
		m.db.Raw(`
			SELECT g.name 
			FROM groups g 
			JOIN activity_groups ag ON g.id = ag.group_id 
			JOIN group_members gm ON g.id = gm.group_id 
			WHERE ag.activity_id = ? AND gm.user_id = ?
		`, activity.ID, userID).Pluck("name", &groupNames)

		activityWithGroups := ActivityWithGroups{
			Activity:   activity,
			GroupNames: groupNames,
		}
		result = append(result, activityWithGroups)
	}

	return result
}

func (m *GormActivityModel) Clear() {
	m.db.Exec("DELETE FROM activity_groups")
	m.db.Exec("DELETE FROM activities")
	m.db.Exec("ALTER SEQUENCE activities_id_seq RESTART WITH 1")
	m.db.Exec("ALTER SEQUENCE activity_groups_id_seq RESTART WITH 1")
}

func (m *GormActivityModel) SeedDefaultData() {
	activity := m.CreateActivity(Activity{
		Title:       "New Activity",
		CreatorID:   1,
		Date:        time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		Description: stringPtr("Sample activity"),
	})
	// Add to group 1
	m.SetActivityGroups(activity.ID, []int{1})
}

func stringPtr(s string) *string { return &s }
