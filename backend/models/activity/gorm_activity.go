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
	return a, true
}

func (m *GormActivityModel) GetActivitiesByCreatorID(creatorID int) []Activity {
	var list []Activity
	m.db.Where("creator_id = ?", creatorID).Find(&list)
	return list
}

func (m *GormActivityModel) GetActivitiesByGroupID(groupID int) []Activity {
	var list []Activity
	m.db.Where("group_id = ?", groupID).Find(&list)
	return list
}

func (m *GormActivityModel) GetActivitiesByGroupIDs(groupIDs []int) []Activity {
	var list []Activity
	m.db.Where("group_id IN ?", groupIDs).Find(&list)
	return list
}

func (m *GormActivityModel) CreateActivity(a Activity) Activity {
	m.db.Create(&a)
	return a
}

func (m *GormActivityModel) UpdateActivity(id int, updates map[string]interface{}) (Activity, bool) {
	var a Activity
	if err := m.db.First(&a, "id = ?", id).Error; err != nil {
		return Activity{}, false
	}
	m.db.Model(&a).Updates(updates)
	return a, true
}

func (m *GormActivityModel) DeleteActivity(id int) bool {
	if err := m.db.Delete(&Activity{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

func (m *GormActivityModel) Clear() {
	m.db.Exec("DELETE FROM activities")
	m.db.Exec("ALTER SEQUENCE activities_id_seq RESTART WITH 1")
}

func (m *GormActivityModel) SeedDefaultData() {
	m.CreateActivity(Activity{
		Title:       "New Activity",
		CreatorID:   1,
		GroupID:     1,
		Date:        time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		Description: stringPtr("Dpzinha legal demais"),
	})
}

func stringPtr(s string) *string { return &s }
