package db

type User struct {
  GUID string
  Username string
  HashedPassword string
}

type Token struct {
  GUID string
  RefreshToken string
}
