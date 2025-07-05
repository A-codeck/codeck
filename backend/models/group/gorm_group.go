package group

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"gorm.io/gorm"
)

type GormGroupModel struct {
	db *gorm.DB
}

func NewGormGroupModel(db *gorm.DB) *GormGroupModel {
	return &GormGroupModel{db: db}
}

func (m *GormGroupModel) GetGroupByID(id int) (Group, bool) {
	var g Group
	if err := m.db.First(&g, "id = ?", id).Error; err != nil {
		return Group{}, false
	}
	return g, true
}

func (m *GormGroupModel) CreateGroup(g Group) Group {
	m.db.Create(&g)
	return g
}

func (m *GormGroupModel) UpdateGroup(id int, updates map[string]interface{}) (Group, bool) {
	var g Group
	if err := m.db.First(&g, "id = ?", id).Error; err != nil {
		return Group{}, false
	}
	m.db.Model(&g).Updates(updates)
	return g, true
}

func (m *GormGroupModel) DeleteGroup(id int) bool {
	if err := m.db.Delete(&Group{}, "id = ?", id).Error; err != nil {
		return false
	}
	return true
}

func (m *GormGroupModel) AddUserToGroup(groupID, userID int) bool {
	member := GroupMember{UserID: userID, GroupID: groupID}
	if err := m.db.Create(&member).Error; err != nil {
		return false
	}
	return true
}

func (m *GormGroupModel) RemoveUserFromGroup(groupID, userID int) bool {
	result := m.db.Delete(&GroupMember{}, "group_id = ? AND user_id = ?", groupID, userID)
	return result.RowsAffected > 0
}

func (m *GormGroupModel) GetGroupMembers(groupID int) ([]GroupMember, bool) {
	var members []GroupMember
	if err := m.db.Where("group_id = ?", groupID).Find(&members).Error; err != nil {
		return nil, false
	}
	return members, true
}

func (m *GormGroupModel) IsUserInGroup(groupID, userID int) bool {
	var member GroupMember
	if err := m.db.Where("group_id = ? AND user_id = ?", groupID, userID).First(&member).Error; err != nil {
		return false
	}
	return true
}

func (m *GormGroupModel) SetUserNickname(groupID, userID int, nickname *string) bool {
	result := m.db.Model(&GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Update("nickname", nickname)
	return result.RowsAffected > 0
}

func (m *GormGroupModel) DeleteUserNickname(groupID, userID int) bool {
	result := m.db.Model(&GroupMember{}).Where("group_id = ? AND user_id = ?", groupID, userID).Update("nickname", nil)
	return result.RowsAffected > 0
}

func generateInviteCode(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "defaultcode" // fallback, should not happen
	}
	return base64.URLEncoding.EncodeToString(b)[:n]
}

func (m *GormGroupModel) CreateInviteLink(groupID, createdBy int, expiresAt *string) (GroupInvite, bool) {
	var expires *time.Time
	if expiresAt != nil {
		t, err := time.Parse(time.RFC3339, *expiresAt)
		if err == nil {
			expires = &t
		}
	}
	inviteCode := generateInviteCode(12)
	invite := GroupInvite{InviteCode: inviteCode, GroupID: groupID, CreatedBy: createdBy, ExpiresAt: expires, IsActive: true}
	if err := m.db.Create(&invite).Error; err != nil {
		return GroupInvite{}, false
	}
	return invite, true
}

func (m *GormGroupModel) GetInviteByCode(inviteCode string) (GroupInvite, bool) {
	var invite GroupInvite
	if err := m.db.Where("invite_code = ?", inviteCode).First(&invite).Error; err != nil {
		return GroupInvite{}, false
	}
	return invite, true
}

func (m *GormGroupModel) DeactivateInvite(inviteCode string) bool {
	result := m.db.Model(&GroupInvite{}).Where("invite_code = ?", inviteCode).Update("is_active", false)
	return result.RowsAffected > 0
}

func (m *GormGroupModel) GetActiveInvites(groupID int) []GroupInvite {
	var invites []GroupInvite
	m.db.Where("group_id = ? AND is_active = ?", groupID, true).Find(&invites)
	return invites
}

func (m *GormGroupModel) GetGroupActivities(groupID int) ([]string, bool) {
	// This would require a join table or relation in a real schema
	return []string{}, true // Placeholder
}

func (m *GormGroupModel) GetUserGroups(userID int) []Group {
	var groups []Group
	m.db.Joins("JOIN group_members ON groups.id = group_members.group_id").
		Where("group_members.user_id = ?", userID).
		Find(&groups)
	return groups
}

func (m *GormGroupModel) Clear() {
	m.db.Exec("DELETE FROM groups")
	m.db.Exec("ALTER SEQUENCE groups_id_seq RESTART WITH 1")
	m.db.Exec("DELETE FROM group_invites")
	m.db.Exec("DELETE FROM group_members")
}

func (m *GormGroupModel) SeedDefaultData() {
	m.CreateGroup(Group{
		ID:          1,
		CreatorID:   1,
		Name:        "New Group",
		StartDate:   time.Now(),
		EndDate:     time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		Description: stringPtr("A test group"),
	})

	m.AddUserToGroup(1, 1)
}

func stringPtr(s string) *string { return &s }
