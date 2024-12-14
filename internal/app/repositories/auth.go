package repositories

import (
  "log"
  "time"

  "golang.org/x/crypto/bcrypt"
  "github.com/golang-jwt/jwt"
)

type AuthRepository interface {
  CreateAccessToken(userID int) string
  CreateRefreshToken(userID int) string
  HashPassword(password string) string
  CompareHashAndPassword(hash string, password string) error
}

type authRepository struct {
  tokenSecret string
}

func (s authRepository) createToken(userID int, scope string) string {
  tokenByte := jwt.New(jwt.SigningMethodHS512)
  now := time.Now().UTC()
  claims := tokenByte.Claims.(jwt.MapClaims)

  claims["sub"] = userID
  claims["exp"] = now.Add(3600).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()
  claims["scope"] = scope

  token, err := tokenByte.SignedString([]byte(s.tokenSecret))
  if err != nil {
    log.Fatal(err.Error())
  }
  return token
}

func (s authRepository) CreateAccessToken(userID int) string {
  token := s.createToken(userID, "access")
  return token
}

func (s authRepository) CreateRefreshToken(userID int) string {
  token := s.createToken(userID, "refresh")
  return token
}

func (s authRepository) HashPassword(password string) string {
  hashedPassword, err := bcrypt.GenerateFromPassword(
    []byte(password), bcrypt.DefaultCost,
  )
  if err != nil {
    log.Fatal(err.Error())
  }
  return string(hashedPassword[:])
}

func (s authRepository) CompareHashAndPassword(hash string, password string) error {
  return bcrypt.CompareHashAndPassword(
    []byte(hash), []byte(password),
  )
}

func NewAuthRepository(tokenSecret string) AuthRepository {
  return &authRepository{tokenSecret: tokenSecret}
}
