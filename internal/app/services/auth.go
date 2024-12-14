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
  Login(schema *schemas.AuthLoginSchema, ip string) (schemas.AuthTokenSchema, error)
  Refresh(user *db.User, token, ip string) (schemas.AuthTokenSchema, error)
  LoginV1(guid, ip string) (schemas.AuthTokenSchema, error)
  RefreshV1(refreshtoken, accessToken, ip string) (schemas.AuthTokenSchema, error)
}

type authService struct {
  authRepository repositories.AuthRepository
  userRepository repositories.UserRepository
  tokenRepository repositories.TokenRepository
  mailRepository repositories.MailRepository
}

func (s authService) LoginV1(guid string, ip string) (schemas.AuthTokenSchema, error) {
  accessToken := s.authRepository.CreateAccessToken(guid, ip)
  refreshToken := s.authRepository.CreateRefreshToken(guid, ip)

  s.tokenRepository.Update(guid, refreshToken)

  return schemas.AuthTokenSchema{
    AccessToken: accessToken,
    RefreshToken: refreshToken,
  }, nil
}

func (s authService) RefreshV1(refreshToken, accessToken, ip string) (schemas.AuthTokenSchema, error) {
  refreshClaims, err := s.authRepository.GetTokenClaims(refreshToken)
  if err != nil {
    log.Print(err)
    return schemas.AuthTokenSchema{}, errors.New("invalid token")
  }
  accessClaims, err := s.authRepository.GetTokenClaims(accessToken)
  if err != nil {
    log.Print(err)
    return schemas.AuthTokenSchema{}, errors.New("invalid token")
  }
  guid := refreshClaims["sub"].(string)

  // if token, err := s.tokenRepository.Get(guid); err != nil || token.RefreshToken != refreshToken {
  if _, err := s.tokenRepository.Get(guid); err != nil {
    log.Print(err)
    return schemas.AuthTokenSchema{}, errors.New("Invalid tokens")
  }

  if refreshClaims["iat"] != accessClaims["iat"] {  // Check if refresh is attached to access
    return schemas.AuthTokenSchema{}, errors.New("Invalid tokens")
  }
  
  tokenIP := refreshClaims["ip"].(string)
  if tokenIP != ip {
    go s.mailRepository.SendIPChangedWarning(guid, ip, tokenIP)
  }

  accessToken = s.authRepository.CreateAccessToken(guid, ip)
  s.tokenRepository.Update(guid, "")

  return schemas.AuthTokenSchema{
    AccessToken: accessToken,
  }, nil
}

func (s authService) Login(schema *schemas.AuthLoginSchema, ip string) (schemas.AuthTokenSchema, error) {
  user, err := s.userRepository.GetByName(schema.Username)
  if err != nil {
    return schemas.AuthTokenSchema{}, errors.New("Invalid username or password")
  }
  // if err := s.authRepository.CompareHashAndPassword(user.HashedPassword, schema.Password); err != nil {
  //  return schemas.AuthTokenSchema{}, errors.New("Invalid username or password")
  // }
  // TODO: Restore after attach database

  accessToken := s.authRepository.CreateAccessToken(user.GUID, ip)
  refreshToken := s.authRepository.CreateRefreshToken(user.GUID, ip)

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
    GUID: model.GUID,
    Username: model.Username,
  }, nil
}

func (s authService) Refresh(user *db.User, token string, ip string) (schemas.AuthTokenSchema, error) {
  tokenClaims, err := s.authRepository.GetTokenClaims(token)
  if err != nil {
    return schemas.AuthTokenSchema{}, err
  }

  tokenIP := tokenClaims["ip"].(string)
  if tokenIP != ip {
    go s.mailRepository.SendIPChangedWarning(user.Username, ip, tokenIP)
  }

  accessToken := s.authRepository.CreateAccessToken(user.GUID, ip)

  return schemas.AuthTokenSchema{
    AccessToken: accessToken,
  }, nil
}

func NewAuthService(
    authRepository repositories.AuthRepository,
    userRepository repositories.UserRepository,
    tokenRepository repositories.TokenRepository,
    mailRepository repositories.MailRepository,
) AuthService {
  return &authService{
    authRepository: authRepository,
    userRepository: userRepository,
    tokenRepository: tokenRepository,
    mailRepository: mailRepository,
  }
}
