package group

type GroupModel interface {
	GetGroupByID(id string) (Group, bool)
	CreateGroup(group Group) Group
	UpdateGroup(id string, updates map[string]interface{}) (Group, bool)
	DeleteGroup(id string) bool
	AddUserToGroup(groupID, userID string) bool
	RemoveUserFromGroup(groupID, userID string) bool
	GetGroupMembers(groupID string) ([]GroupMember, bool)
	IsUserInGroup(groupID, userID string) bool
	SetUserNickname(groupID, userID string, nickname *string) bool
	DeleteUserNickname(groupID, userID string) bool
	CreateInviteLink(groupID, createdBy string, expiresAt *string) (GroupInvite, bool)
	GetInviteByCode(inviteCode string) (GroupInvite, bool)
	DeactivateInvite(inviteCode string) bool
	GetActiveInvites(groupID string) []GroupInvite
	GetGroupActivities(groupID string) ([]string, bool)
	GetUserGroups(userID string) []Group
}

var DefaultGroupModel GroupModel = NewInMemoryGroup()
