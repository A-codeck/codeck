package group

type GroupModel interface {
	GetGroupByID(id string) (Group, bool)
	CreateGroup(group Group) Group
	UpdateGroup(id string, updates map[string]interface{}) (Group, bool)
	DeleteGroup(id string) bool
	AddUserToGroup(groupID, userID string, nickname *string) bool
	RemoveUserFromGroup(groupID, userID string) bool
	GetGroupMembers(groupID string) ([]GroupMember, bool)
	IsUserInGroup(groupID, userID string) bool
	CreateInviteLink(groupID, createdBy string, expiresAt *string) (GroupInvite, bool)
	GetInviteByCode(inviteCode string) (GroupInvite, bool)
	DeactivateInvite(inviteCode string) bool
	GetActiveInvites(groupID string) []GroupInvite
}

var DefaultGroupModel GroupModel = NewInMemoryGroup()
