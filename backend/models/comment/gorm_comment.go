package comment

import "gorm.io/gorm"

type GormCommentModel struct {
	db *gorm.DB
}

func NewGormCommentModel(db *gorm.DB) *GormCommentModel {
	return &GormCommentModel{db: db}
}

func (m *GormCommentModel) GetCommentByID(id int) (Comment, bool) {
	var c Comment
	if err := m.db.First(&c, "id = ?", id).Error; err != nil {
		return Comment{}, false
	}
	return c, true
}

func (m *GormCommentModel) GetCommentsByActivityID(activityID int) []Comment {
	var list []Comment
	m.db.Where("activity_id = ?", activityID).Find(&list)
	return list
}

func (m *GormCommentModel) CreateComment(c Comment) Comment {
	m.db.Create(&c)
	return c
}

func (m *GormCommentModel) DeleteComment(id int) bool {
	if err := m.db.Delete(&Comment{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

func (m *GormCommentModel) Clear() {
	m.db.Exec("DELETE FROM comments")
	m.db.Exec("ALTER SEQUENCE comments_id_seq RESTART WITH 1")
}

func (m *GormCommentModel) SeedDefaultData() {
	m.CreateComment(Comment{
		ActivityID: 1,
		UserID:     1,
		Content:    "Great activity!",
	})
}
