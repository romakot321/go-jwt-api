package repositories

import (
  "github.com/romakot321/go-jwt-api/internal/app/db"
)

type UserRepository interface {
  Create(model *db.User) error
  Get(modelID string) (*db.User, error)
  GetByName(name string) (*db.User, error)
}

type userRepository struct {

}

func (s userRepository) Create(model *db.User) error {
  model.GUID = "guid"
  return nil
}

func (s userRepository) Get(modelID string) (*db.User, error) {
  return &db.User{GUID: modelID, Username: "test", HashedPassword: "$2a$10$2nLHvswEFnZVQe.OeJ2eAe6nwnKYDI84bQGJoEFm735mp/iS0ecZG"}, nil
}

func (s userRepository) GetByName(name string) (*db.User, error) {
  return &db.User{GUID: "guid", Username: name, HashedPassword: "$2a$10$2nLHvswEFnZVQe.OeJ2eAe6nwnKYDI84bQGJoEFm735mp/iS0ecZG"}, nil
}

func NewUserRepository() UserRepository {
  return &userRepository{}
}

func GetUserByID(userID string) (*db.User, error) {
  return &db.User{GUID: userID, Username: "test", HashedPassword: "$2a$10$2nLHvswEFnZVQe.OeJ2eAe6nwnKYDI84bQGJoEFm735mp/iS0ecZG"}, nil
}
