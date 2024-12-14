package middleware

import (
  "strings"
  "errors"
  "log"

  "github.com/gofiber/fiber/v2"
  "github.com/golang-jwt/jwt"
  "github.com/romakot321/go-jwt-api/internal/app/repositories"
  "github.com/romakot321/go-jwt-api/internal/app/db"
)

func authenticateUser(c *fiber.Ctx, scope string) error {
  token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
  if token == "" {
    token = c.Cookies(scope + "-token")
  }
  tokenByte, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
    if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, errors.New("Invalid token signing method")
    }
    return []byte("REPLACEME"), nil
  })
  if err != nil {
    log.Print(err, tokenByte)
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }

  claims, ok := tokenByte.Claims.(jwt.MapClaims)
  if !ok || !tokenByte.Valid {
    log.Print("err2")
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }
  tokenScope, ok := claims["scope"]
  if !ok || !strings.Contains(tokenScope.(string), scope) {
    log.Print("err")
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }

  userModel, err := repositories.GetUserByID(claims["sub"].(string))
  if err != nil {
    log.Print("err")
    return c.Status(401).JSON(fiber.Map{
      "status": "fail",
      "message": "invalid token",
    })
  }

  c.Locals("user", userModel)
  return nil
}

func AuthenticateUserAccess(c *fiber.Ctx) error {
  authenticateUser(c, "access")
  if _, ok := c.Locals("user").(*db.User); ok {
    return c.Next()
  }
  return nil
}

func AuthenticateUserRefresh(c *fiber.Ctx) error {
  authenticateUser(c, "refresh")
  if _, ok := c.Locals("user").(*db.User); ok {
    return c.Next()
  }
  return nil
}
