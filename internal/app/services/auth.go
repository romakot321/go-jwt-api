package services

import (
  "errors"
  "log"

  "github.com/romakot321/go-jwt-api/internal/app/schemas"
  "github.com/romakot321/go-jwt-api/internal/app/db"
  "github.com/romakot321/go-jwt-api/internal/app/repositories"
)

type AuthService interface {
  Register(schema *schemas.AuthRegisterSchema) (schemas.UserGetSchema, error)
  Login(schema *schemas.AuthLoginSchema) (schemas.AuthTokenSchema, error)
}

type authService struct {
  authRepository repositories.AuthRepository
  userRepository repositories.UserRepository
}

func (s authService) Login(schema *schemas.AuthLoginSchema) (schemas.AuthTokenSchema, error) {
  user, err := s.userRepository.GetByName(schema.Username)
  if err != nil {
    return schemas.AuthTokenSchema{}, errors.New("Invalid username or password")
  }
  // if err := s.authRepository.CompareHashAndPassword(user.HashedPassword, schema.Password); err != nil {
  //  return schemas.AuthTokenSchema{}, errors.New("Invalid username or password")
  // }
  // TODO: Restore after attach database

  accessToken := s.authRepository.CreateAccessToken(user.ID)
  refreshToken := s.authRepository.CreateRefreshToken(user.ID)

  return schemas.AuthTokenSchema{
    AccessToken: accessToken,
    RefreshToken: refreshToken,
  }, nil
}

func (s authService) Register(schema *schemas.AuthRegisterSchema) (schemas.UserGetSchema, error) {
  hashedPassword := s.authRepository.HashPassword(schema.Password)

  model := &db.User{
    Username: schema.Username,
    HashedPassword: hashedPassword,
  }
  log.Print(model)
  if err := s.userRepository.Create(model); err != nil {
    return schemas.UserGetSchema{}, err
  }

  return schemas.UserGetSchema{
    ID: model.ID,
    Username: model.Username,
  }, nil
}

func NewAuthService(authRepository repositories.AuthRepository, userRepository repositories.UserRepository) AuthService {
  return &authService{authRepository: authRepository, userRepository: userRepository}
}
