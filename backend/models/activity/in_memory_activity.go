package activity

import (
	"strconv"
	"sync"
	"time"
)

type inMemoryActivity struct {
	activitys map[string]Activity
	idCounter int
	mutex     sync.Mutex
}

type InMemoryActivityModel = inMemoryActivity

func NewInMemoryActivity() *inMemoryActivity {
	return &inMemoryActivity{
		activitys: make(map[string]Activity),
		idCounter: 1,
	}
}

func (g *inMemoryActivity) GetActivityByID(id string) (Activity, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	activity, exists := g.activitys[id]
	return activity, exists
}

func (g *inMemoryActivity) CreateActivity(activity Activity) Activity {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	activity.ID = strconv.Itoa(g.idCounter)
	g.idCounter++
	g.activitys[activity.ID] = activity
	return activity
}

func (g *inMemoryActivity) UpdateActivity(id string, updates map[string]interface{}) (Activity, bool) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	activity, exists := g.activitys[id]
	if !exists {
		return Activity{}, false
	}

	if description, ok := updates["description"].(string); ok {
		activity.Description = &description
	}
	if activityImage, ok := updates["activity_image"].(string); ok {
		activity.ActivityImage = &activityImage
	}

	g.activitys[id] = activity
	return activity, true
}

func (g *inMemoryActivity) DeleteActivity(id string) bool {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	_, exists := g.activitys[id]
	if exists {
		delete(g.activitys, id)
	}
	return exists
}

func (g *inMemoryActivity) Clear() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.activitys = make(map[string]Activity)
	g.idCounter = 1
}

func (g *inMemoryActivity) SeedDefaultData() {
	date := time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	g.CreateActivity(Activity{
		CreatorID: "1",
		Title:     "Default Activity",
		Date:      date,
	})
}
