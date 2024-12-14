package db

type User struct {
  GUID string `gorm:"type:varchar(36);primary_key"`
  Username string `gorm:"type:varchar(100);not null"`
  HashedPassword string `gorm:"type:varchar(100);not null"`
}

type Token struct {
  GUID string `gorm:"type:varchar(36);primary_key"`
  RefreshToken string `gorm:"type:varchar(200);not null"`
}
