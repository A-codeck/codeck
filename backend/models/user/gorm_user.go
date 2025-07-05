package user

import "gorm.io/gorm"

type GormUserModel struct {
	db *gorm.DB
}

func NewGormUserModel(db *gorm.DB) *GormUserModel {
	return &GormUserModel{db: db}
}

func (m *GormUserModel) GetUserByID(id int) (User, bool) {
	var u User
	if err := m.db.First(&u, "id = ?", id).Error; err != nil {
		return User{}, false
	}
	return u, true
}

func (m *GormUserModel) GetUserByEmail(email string) (User, bool) {
	var u User
	if err := m.db.First(&u, "email = ?", email).Error; err != nil {
		return User{}, false
	}
	return u, true
}

func (m *GormUserModel) CreateUser(u User) User {
	m.db.Create(&u)
	return u
}

func (m *GormUserModel) ValidateCredentials(email, password string) (User, bool) {
	var u User
	if err := m.db.First(&u, "email = ?", email).Error; err != nil {
		return User{}, false
	}
	if u.Password != password {
		return User{}, false
	}
	return u, true
}

func (m *GormUserModel) Clear() {
	m.db.Exec("DELETE FROM users")
	m.db.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func (m *GormUserModel) SeedDefaultData() {
	m.CreateUser(User{
		ID:       1,
		Email:    "user@example.com",
		Name:     "Test User",
		Password: "password123",
	})
	m.CreateUser(User{
		ID:       2,
		Email:    "user2@example.com",
		Name:     "Test User 2",
		Password: "password123",
	})
}
