package repositories

import (
  "github.com/romakot321/go-jwt-api/internal/app/db"
)

type TokenRepository interface {
  Create(model *db.Token) error
  Get(modelID string) (*db.Token, error)
  Update(guid string, refreshToken string) error
}

type tokenRepository struct {

}

func (s tokenRepository) Create(model *db.Token) error {
  return nil
}

func (s tokenRepository) Get(guid string) (*db.Token, error) {
  return &db.Token{GUID: guid, RefreshToken: "refresh"}, nil
}

func (s tokenRepository) Update(guid string, refreshToken string) error {
  return nil
}

func NewTokenRepository() TokenRepository {
  return &tokenRepository{}
}
