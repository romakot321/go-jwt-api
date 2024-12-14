package repositories

import (
  "github.com/romakot321/go-jwt-api/internal/app/db"
)

type UserRepository interface {
  Create(model *db.User) error
  Get(modelID int) (*db.User, error)
  GetByName(name string) (*db.User, error)
}

type userRepository struct {

}

func (s userRepository) Create(model *db.User) error {
  model.ID = 1
  return nil
}

func (s userRepository) Get(modelID int) (*db.User, error) {
  return &db.User{ID: modelID, Username: "test", HashedPassword: "$2a$10$2nLHvswEFnZVQe.OeJ2eAe6nwnKYDI84bQGJoEFm735mp/iS0ecZG"}, nil
}

func (s userRepository) GetByName(name string) (*db.User, error) {
  return &db.User{ID: 1, Username: name, HashedPassword: "$2a$10$2nLHvswEFnZVQe.OeJ2eAe6nwnKYDI84bQGJoEFm735mp/iS0ecZG"}, nil
}

func NewUserRepository() UserRepository {
  return &userRepository{}
}

func GetUserByID(userID int) (*db.User, error) {
  return &db.User{ID: userID, Username: "test", HashedPassword: "$2a$10$2nLHvswEFnZVQe.OeJ2eAe6nwnKYDI84bQGJoEFm735mp/iS0ecZG"}, nil
}
