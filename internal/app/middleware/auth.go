package middleware

import (
  "strings"
  "errors"
  "log"

  "github.com/gofiber/fiber/v2"
  "github.com/golang-jwt/jwt"
  "github.com/romakot321/go-jwt-api/internal/app/repositories"
)

func authenticateUser(c *fiber.Ctx, scope string) error {
  token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
  if token == "" {
    token = c.Cookies("access-token")
  }
  tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
    if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("Invalid token signing method")
    }
    return []byte("REPLACEME"), nil
  })
  if err != nil {
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }

  claims, ok := tokenByte.Claims.(jwt.MapClaims)
  if !ok || !tokenByte.Valid {
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }
  tokenScope, ok := claims["scope"]
  if !ok || !strings.Contains(tokenScope.(string), scope) {
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }

  userModel, err := repositories.GetUserByID(claims["sub"].(int))
  if err != nil {
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }

  c.Locals("user", userModel)
  return nil
}

func AuthenticateUserAccess(c *fiber.Ctx) error {
  if err := authenticateUser(c, "access"); err != nil {
    log.Print("err")
    return err
  }
  return c.Next()
}

func AuthenticateUserRefresh(c *fiber.Ctx) error {
  if err := authenticateUser(c, "refresh"); err != nil {
    return err
  }
  return c.Next()
}
