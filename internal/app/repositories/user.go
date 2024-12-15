package repositories

import (
  "errors"
  "strings"

  "github.com/romakot321/go-jwt-api/internal/app/db"
)

type UserRepository interface {
  Create(model *db.User) error
  Get(modelID string) (*db.User, error)
  GetByName(name string) (*db.User, error)
}

type userRepository struct {
  // TODO: Add db connection as dependency
}

func (s userRepository) Create(model *db.User) error {
  result := db.DB.Create(model)
  if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return errors.New("User with this username already exists")
	} else if result.Error != nil {
		return errors.New("Bad gateway")
	}
  return nil
}

func (s userRepository) Get(modelID string) (*db.User, error) {
  var user db.User
  result := db.DB.First(&user, "guid = ?", modelID)
	if result.Error != nil {
		return &user, errors.New("User not found")
	}
  return& user, nil
}

func (s userRepository) GetByName(name string) (*db.User, error) {
  var user db.User
  result := db.DB.First(&user, "username = ?", name)
	if result.Error != nil {
		return &user, errors.New("User not found")
	}
  return &user, nil
}

func NewUserRepository() UserRepository {
  return &userRepository{}
}

func GetUserByID(userID string) (*db.User, error) {
  var user db.User
  result := db.DB.First(&user, "guid = ?", userID)
	if result.Error != nil {
		return &user, errors.New("User not found")
	}
  return &user, nil
}
