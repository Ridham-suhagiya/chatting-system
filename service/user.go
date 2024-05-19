package service

import (
	"chatting-system-backend/database"
	"chatting-system-backend/model"
	"fmt"
)

// User struct defines the structure of the user entit

// UserRepository interface defines the methods to be implemented
type UserService interface {
	CreateUser(user model.User) error
	GetUserByID(userID string) (*model.User, error)
	GetAllUsers() (*[]model.User, error)
}

// userService struct holds a reference to the database connection
type userService struct {
	DB *database.DB
}

// NewUserService returns a new instance of UserService
func NewUserService(db *database.DB) UserService {
	return &userService{DB: db}
}

// CreateUser method creates a new user in the database
func (s *userService) CreateUser(user model.User) error {
	fmt.Printf("Creating user: %+v\n", user)
	return s.DB.Create(&user).Error
}

// GetUserByID method retrieves a user by ID from the database
func (s *userService) GetUserByID(userID string) (*model.User, error) {
	var user model.User
	if err := s.DB.First(&user, `id=?`, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *userService) GetAllUsers() (*[]model.User, error) {
	var users []model.User
	if err := s.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
