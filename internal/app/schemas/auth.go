package schemas

type AuthRegisterSchema struct {
  GUID string `json:"guid"`
  Username string `json:"username"`
  Password string `json:"password"`
}

type AuthLoginSchema struct {
  Username string `json:"username"`
  Password string `json:"password"`
}

type AuthTokenSchema struct {
  AccessToken string `json:"access_token"`
  RefreshToken string `json:"refresh_token"`
}
