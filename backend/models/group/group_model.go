package group

type GroupModel interface {
	GetGroupByID(id string) (Group, bool)
	CreateGroup(group Group) Group
	UpdateGroup(id string, updates map[string]interface{}) (Group, bool)
	DeleteGroup(id string) bool
}

var DefaultGroupModel GroupModel = NewInMemoryGroup()
