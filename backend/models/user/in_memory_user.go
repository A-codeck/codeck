package user

import (
	"strconv"
	"strings"
	"sync"
)

type inMemoryUser struct {
	users     map[string]User
	emailMap  map[string]string 
	idCounter int
	mutex     sync.Mutex
}

type InMemoryUserModel = inMemoryUser

func NewInMemoryUser() *inMemoryUser {
	return &inMemoryUser{
		users:     make(map[string]User),
		emailMap:  make(map[string]string),
		idCounter: 1,
	}
}

func (u *inMemoryUser) GetUserByID(id string) (User, bool) {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	user, exists := u.users[id]
	return user, exists
}

func (u *inMemoryUser) GetUserByEmail(email string) (User, bool) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	email = strings.ToLower(email)
	id, exists := u.emailMap[email]
	if !exists {
		return User{}, false
	}

	user, exists := u.users[id]
	return user, exists
}

func (u *inMemoryUser) CreateUser(user User) User {
	u.mutex.Lock()
	defer u.mutex.Unlock()

 	user.ID = strconv.Itoa(u.idCounter)
	u.idCounter++

	u.users[user.ID] = user

	u.emailMap[strings.ToLower(user.Email)] = user.ID

	return user
}

func (u *inMemoryUser) ValidateCredentials(email, password string) (User, bool) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	email = strings.ToLower(email)
	id, exists := u.emailMap[email]
	if !exists {
		return User{}, false
	}

	user, exists := u.users[id]
	if !exists {
		return User{}, false
	}

	// Simple password check
	if user.Password != password {
		return User{}, false
	}

	return user, true
}

func (u *inMemoryUser) Clear() {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	u.users = make(map[string]User)
	u.emailMap = make(map[string]string)
	u.idCounter = 1
}

func (u *inMemoryUser) SeedDefaultData() {
	// Create a test user
	testUser := User{
		Email:    "user@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	u.CreateUser(testUser)
}
