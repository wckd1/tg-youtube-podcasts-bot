package user

type UserRepository interface {
	CreateUser(*User) error
	GetUser(string) (User, error)
	DeleteUser(string) error
}
