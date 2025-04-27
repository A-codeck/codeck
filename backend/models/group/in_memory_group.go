package group

import (
	"strconv"
	"sync"
	"time"
)

type inMemoryGroup struct {
	groups    map[string]Group
	idCounter int
	mutex     sync.Mutex
}

type InMemoryGroupModel = inMemoryGroup

func NewInMemoryGroup() *inMemoryGroup {
	return &inMemoryGroup{
		groups:    make(map[string]Group),
		idCounter: 1,
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
