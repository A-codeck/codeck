package group

import (
	"strconv"
	"sync"
	"time"
)

type inMemoryGroup struct {
	groups          map[string]Group
	groupMembers    map[string][]GroupMember // groupID -> []GroupMember
	invites         map[string]GroupInvite   // inviteCode -> GroupInvite
	groupActivities map[string][]string      // groupID -> []activityID
	idCounter       int
	mutex           sync.Mutex
}

type InMemoryGroupModel = inMemoryGroup

func NewInMemoryGroup() *inMemoryGroup {
	return &inMemoryGroup{
		groups:          make(map[string]Group),
		groupMembers:    make(map[string][]GroupMember),
		invites:         make(map[string]GroupInvite),
		groupActivities: make(map[string][]string),
		idCounter:       1,
	}
}

func (g *inMemoryGroup) GetGroupByID(id string) (Group, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	group, exists := g.groups[id]
	return group, exists
}

func (g *inMemoryGroup) CreateGroup(group Group) Group {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	group.ID = strconv.Itoa(g.idCounter)
	g.idCounter++
	g.groups[group.ID] = group

	creatorMember := GroupMember{
		UserID:   group.CreatorID,
		GroupID:  group.ID,
		Nickname: nil,
	}
	g.groupMembers[group.ID] = append(g.groupMembers[group.ID], creatorMember)

	return group
}

func (g *inMemoryGroup) UpdateGroup(id string, updates map[string]interface{}) (Group, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	group, exists := g.groups[id]
	if !exists {
		return Group{}, false
	}

	if description, ok := updates["description"].(string); ok {
		group.Description = &description
	}
	if groupImage, ok := updates["group_image"].(string); ok {
		group.GroupImage = &groupImage
	}
	if endDate, ok := updates["end_date"].(string); ok {
		group.EndDate = endDate
	}

	g.groups[id] = group
	return group, true
}

func (g *inMemoryGroup) DeleteGroup(id string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	_, exists := g.groups[id]
	if exists {
		delete(g.groups, id)
	}
	return exists
}

func (g *inMemoryGroup) Clear() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.groups = make(map[string]Group)
	g.groupMembers = make(map[string][]GroupMember)
	g.invites = make(map[string]GroupInvite)
	g.groupActivities = make(map[string][]string)
	g.idCounter = 1
}

func (g *inMemoryGroup) SeedDefaultData() {
	start := time.Now().Format("2006-01-02")
	end := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	g.CreateGroup(Group{
		CreatorID: "1",
		Name:      "Default Group",
		StartDate: start,
		EndDate:   end,
	})
}

func (g *inMemoryGroup) AddUserToGroup(groupID, userID string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.groups[groupID]; !exists {
		return false
	}

	if members, exists := g.groupMembers[groupID]; exists {
		for _, member := range members {
			if member.UserID == userID {
				return false
			}
		}
	}

	newMember := GroupMember{
		UserID:   userID,
		GroupID:  groupID,
		Nickname: nil,
	}

	g.groupMembers[groupID] = append(g.groupMembers[groupID], newMember)
	return true
}

func (g *inMemoryGroup) RemoveUserFromGroup(groupID, userID string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	members, exists := g.groupMembers[groupID]
	if !exists {
		return false
	}

	for i, member := range members {
		if member.UserID == userID {
			g.groupMembers[groupID] = append(members[:i], members[i+1:]...)
			return true
		}
	}

	return false
}

func (g *inMemoryGroup) GetGroupMembers(groupID string) ([]GroupMember, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.groups[groupID]; !exists {
		return nil, false
	}

	members, exists := g.groupMembers[groupID]
	if !exists {
		return []GroupMember{}, true
	}

	return members, true
}

func (g *inMemoryGroup) IsUserInGroup(groupID, userID string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	members, exists := g.groupMembers[groupID]
	if !exists {
		return false
	}

	for _, member := range members {
		if member.UserID == userID {
			return true
		}
	}

	return false
}

func (g *inMemoryGroup) SetUserNickname(groupID, userID string, nickname *string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	members, exists := g.groupMembers[groupID]
	if !exists {
		return false
	}

	for i, member := range members {
		if member.UserID == userID {
			members[i].Nickname = nickname
			g.groupMembers[groupID] = members
			return true
		}
	}

	return false
}

func (g *inMemoryGroup) DeleteUserNickname(groupID, userID string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	members, exists := g.groupMembers[groupID]
	if !exists {
		return false
	}

	for i, member := range members {
		if member.UserID == userID {
			members[i].Nickname = nil
			g.groupMembers[groupID] = members
			return true
		}
	}

	return false
}

func (g *inMemoryGroup) CreateInviteLink(groupID, createdBy string, expiresAt *string) (GroupInvite, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.groups[groupID]; !exists {
		return GroupInvite{}, false
	}

	inviteCode := g.generateInviteCode()

	invite := GroupInvite{
		InviteCode: inviteCode,
		GroupID:    groupID,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now().Format("2006-01-02T15:04:05Z"),
		ExpiresAt:  expiresAt,
		IsActive:   true,
	}

	g.invites[inviteCode] = invite
	return invite, true
}

func (g *inMemoryGroup) GetInviteByCode(inviteCode string) (GroupInvite, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	invite, exists := g.invites[inviteCode]
	if !exists {
		return GroupInvite{}, false
	}

	if invite.ExpiresAt != nil {
		expiryTime, err := time.Parse("2006-01-02T15:04:05Z", *invite.ExpiresAt)
		if err == nil && time.Now().After(expiryTime) {
			invite.IsActive = false
			g.invites[inviteCode] = invite
			return invite, false
		}
	}

	return invite, invite.IsActive
}

func (g *inMemoryGroup) DeactivateInvite(inviteCode string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	invite, exists := g.invites[inviteCode]
	if !exists {
		return false
	}

	invite.IsActive = false
	g.invites[inviteCode] = invite
	return true
}

func (g *inMemoryGroup) GetActiveInvites(groupID string) []GroupInvite {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var activeInvites []GroupInvite

	for _, invite := range g.invites {
		if invite.GroupID == groupID && invite.IsActive {
			if invite.ExpiresAt != nil {
				expiryTime, err := time.Parse("2006-01-02T15:04:05Z", *invite.ExpiresAt)
				if err == nil && time.Now().After(expiryTime) {
					invite.IsActive = false
					g.invites[invite.InviteCode] = invite
					continue
				}
			}
			activeInvites = append(activeInvites, invite)
		}
	}

	return activeInvites
}

func (g *inMemoryGroup) GetGroupActivities(groupID string) ([]string, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if _, exists := g.groups[groupID]; !exists {
		return nil, false
	}

	activities, exists := g.groupActivities[groupID]
	if !exists {
		return []string{}, true
	}

	return activities, true
}

func (g *inMemoryGroup) GetUserGroups(userID string) []Group {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var userGroups []Group
	for groupID, members := range g.groupMembers {
		for _, member := range members {
			if member.UserID == userID {
				if group, exists := g.groups[groupID]; exists {
					userGroups = append(userGroups, group)
				}
				break
			}
		}
	}

	return userGroups
}

func (g *inMemoryGroup) generateInviteCode() string {
	bytes := make([]byte, 4)
	now := time.Now().UnixNano()
	for i := range bytes {
		bytes[i] = byte((now >> (i * 8)) % 256)
	}

	code := ""
	hexChars := "0123456789ABCDEF"
	for _, b := range bytes {
		code += string(hexChars[b>>4]) + string(hexChars[b&0xF])
	}

	if _, exists := g.invites[code]; exists {
		return g.generateInviteCode()
	}

	return code
}
