package user

type UserModel interface {
	GetUserByID(id int) (User, bool)
	GetUserByEmail(email string) (User, bool)
	CreateUser(user User) User
	ValidateCredentials(email, password string) (User, bool)
}

// DefaultUserModel must be set in main.go after DB initialization
var DefaultUserModel UserModel
