package repositories

import (
  "strings"
  "errors"

  "github.com/romakot321/go-jwt-api/internal/app/db"
)

type TokenRepository interface {
  Create(model *db.Token) error
  Get(modelID string) (*db.Token, error)
  Update(guid string, refreshToken string) error
  UpdateOrCreate(guid string, refreshToken string) error
}

type tokenRepository struct {
  // TODO: Add db connection as dependency
}

func (s tokenRepository) Create(model *db.Token) error {
  result := db.DB.Create(model)
  if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return errors.New("Token with this guid already exists")
	} else if result.Error != nil {
		return errors.New("Bad gateway")
	}
  return nil}

func (s tokenRepository) Get(guid string) (*db.Token, error) {
  var model db.Token
  result := db.DB.First(&model, "guid = ?", guid)
	if result.Error != nil {
		return &model, errors.New("Model not found")
	}
  return &model, nil
}

func (s tokenRepository) Update(guid string, refreshToken string) error {
  var model db.Token
  if err := db.DB.Model(&model).Where("guid = ?", guid).Update("refresh_token", refreshToken).Error; err != nil {
    return err
  }
  return nil
}

func (s tokenRepository) UpdateOrCreate(guid string, refreshToken string) error {
  model := db.Token{GUID: guid, RefreshToken: refreshToken}
  result := db.DB.Model(&model).Where("guid = ?", guid).Update("refresh_token", refreshToken)
  if result.Error != nil {
    return result.Error
  }
  if result.RowsAffected == 0 {
    db.DB.Create(&model)
  }
  return nil
}

func NewTokenRepository() TokenRepository {
  return &tokenRepository{}
}
