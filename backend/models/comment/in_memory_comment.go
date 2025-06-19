package comment

import (
	"strconv"
	"sync"
	"time"
)

type inMemoryComment struct {
	comments  map[string]Comment
	idCounter int
	mutex     sync.Mutex
}

type InMemoryCommentModel = inMemoryComment

func NewInMemoryComment() *inMemoryComment {
	return &inMemoryComment{
		comments:  make(map[string]Comment),
		idCounter: 1,
	}
}

func (c *inMemoryComment) GetCommentByID(id string) (Comment, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	comment, exists := c.comments[id]
	return comment, exists
}

func (c *inMemoryComment) GetCommentsByActivityID(activityID string) []Comment {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var comments []Comment
	for _, comment := range c.comments {
		if comment.ActivityID == activityID {
			comments = append(comments, comment)
		}
	}
	return comments
}

func (c *inMemoryComment) CreateComment(comment Comment) Comment {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	comment.ID = strconv.Itoa(c.idCounter)
	c.idCounter++
	comment.CreatedAt = time.Now().Format("2006-01-02T15:04:05Z")

	c.comments[comment.ID] = comment
	return comment
}

func (c *inMemoryComment) DeleteComment(id string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	_, exists := c.comments[id]
	if exists {
		delete(c.comments, id)
	}
	return exists
}

func (c *inMemoryComment) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.comments = make(map[string]Comment)
	c.idCounter = 1
}

func (c *inMemoryComment) SeedDefaultData() {
	c.CreateComment(Comment{
		ActivityID: "1",
		UserID:     "1",
		Content:    "Great activity!",
	})
}
