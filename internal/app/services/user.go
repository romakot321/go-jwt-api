package services

import (
  "github.com/romakot321/go-jwt-api/internal/app/repositories"
)

type UserService interface {
}

type userService struct {
  userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
  return &userService{userRepository: userRepository}
}
