package repositories

import (
  "log"
  "time"
  "errors"

  "golang.org/x/crypto/bcrypt"
  "github.com/golang-jwt/jwt"
)

type AuthRepository interface {
  CreateAccessToken(userID, ip string) string
  CreateRefreshToken(userID, ip string) string
  HashPassword(password string) string
  CompareHashAndPassword(hash string, password string) error
  GetTokenClaims(token string) (jwt.MapClaims, error)
}

type authRepository struct {
  tokenSecret string
}

func (s authRepository) createToken(userID string, scope string, exp int, ip string) string {
  tokenByte := jwt.New(jwt.SigningMethodHS512)
  now := time.Now().UTC()
  claims := tokenByte.Claims.(jwt.MapClaims)

  claims["sub"] = userID
  claims["exp"] = now.Add(time.Duration(exp)).Unix()
  claims["iat"] = now.Unix()
  claims["nbf"] = now.Unix()
  claims["scope"] = scope
  claims["ip"] = ip

  token, err := tokenByte.SignedString([]byte(s.tokenSecret))
  if err != nil {
    log.Fatal(err.Error())
  }
  return token
}

func (s authRepository) CreateAccessToken(userID string, ip string) string {
  return s.createToken(userID, "access", 60 * 60 * 1000000, ip)
}

func (s authRepository) CreateRefreshToken(userID string, ip string) string {
  return s.createToken(userID, "refresh", 48 * 60 * 60 * 1000000, ip)
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

func (s authRepository) GetTokenClaims(token string) (jwt.MapClaims, error) {
  tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
    if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("Invalid token signing method")
    }
    return []byte(s.tokenSecret), nil
  })
  if err != nil {
    log.Print(err)
    return jwt.MapClaims{}, errors.New("Invalid token")
  }

  claims, ok := tokenByte.Claims.(jwt.MapClaims)
  if !ok || !tokenByte.Valid {
    return jwt.MapClaims{}, errors.New("Invalid token")
  }
  return claims, nil
}

func NewAuthRepository(tokenSecret string) AuthRepository {
  return &authRepository{tokenSecret: tokenSecret}
}
