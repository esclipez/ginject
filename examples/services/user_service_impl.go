package services

import (
	"context"
	"fmt"

	"github.com/esclipez/ginject/examples/interfaces"
)

// UserServiceImpl implements UserService interface
type UserServiceImpl struct {
	initialized bool
}

func NewUserService() *UserServiceImpl {
	return &UserServiceImpl{}
}

func (s *UserServiceImpl) Init(ctx context.Context) error {
	fmt.Println("[UserService] Initializing...")
	s.initialized = true
	return nil
}

func (s *UserServiceImpl) Name() string {
	return "UserSrv"
}

func (s *UserServiceImpl) Start(ctx context.Context) error {
	fmt.Println("[UserService] Starting...")
	return nil
}

func (s *UserServiceImpl) Stop(ctx context.Context) error {
	fmt.Println("[UserService] Stopping...")
	return nil
}

func (s *UserServiceImpl) GetUser(id string) (*interfaces.User, error) {
	if !s.initialized {
		return nil, fmt.Errorf("service not initialized")
	}
	return &interfaces.User{
		ID:    id,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil
}

func (s *UserServiceImpl) CreateUser(user *interfaces.User) error {
	if !s.initialized {
		return fmt.Errorf("service not initialized")
	}
	fmt.Printf("Creating user: %+v\n", user)
	return nil
}
