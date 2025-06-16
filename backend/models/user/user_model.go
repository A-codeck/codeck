package user

type UserModel interface {
	GetUserByID(id string) (User, bool)
	GetUserByEmail(email string) (User, bool)
	CreateUser(user User) User
	ValidateCredentials(email, password string) (User, bool)
}

var DefaultUserModel UserModel = NewInMemoryUser()
