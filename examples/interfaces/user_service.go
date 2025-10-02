package interfaces

// UserService interface for user operations
type UserService interface {
	GetUser(id string) (*User, error)
	CreateUser(user *User) error
}

// User represents a user entity
type User struct {
	ID    string
	Name  string
	Email string
}
