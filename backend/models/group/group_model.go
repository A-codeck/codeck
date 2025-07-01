package group

type GroupModel interface {
	GetGroupByID(id int) (Group, bool)
	CreateGroup(group Group) Group
	UpdateGroup(id int, updates map[string]interface{}) (Group, bool)
	DeleteGroup(id int) bool
	AddUserToGroup(groupID, userID int) bool
	RemoveUserFromGroup(groupID, userID int) bool
	GetGroupMembers(groupID int) ([]GroupMember, bool)
	IsUserInGroup(groupID, userID int) bool
	SetUserNickname(groupID, userID int, nickname *string) bool
	DeleteUserNickname(groupID, userID int) bool
	CreateInviteLink(groupID, createdBy int, expiresAt *string) (GroupInvite, bool)
	GetInviteByCode(inviteCode string) (GroupInvite, bool)
	DeactivateInvite(inviteCode string) bool
	GetActiveInvites(groupID int) []GroupInvite
	GetGroupActivities(groupID int) ([]string, bool)
}

// DefaultGroupModel must be set in main.go after DB initialization
var DefaultGroupModel GroupModel
