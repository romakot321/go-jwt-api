package controllers

import (
  "strings"

  "github.com/gofiber/fiber/v2"
  "github.com/romakot321/go-jwt-api/internal/app/schemas"
  "github.com/romakot321/go-jwt-api/internal/app/services"
  "github.com/romakot321/go-jwt-api/internal/app/middleware"
  "github.com/romakot321/go-jwt-api/internal/app/db"
)

type AuthController interface {
  Register(router fiber.Router)
}

type authController struct {
  authService services.AuthService
}

func (c authController) Register(router fiber.Router) {
  router.Post("/register", c.register)
  router.Post("/login", c.login)
  router.Post("/refresh", middleware.AuthenticateUserRefresh, c.refresh)

  router.Post("/v1/login", c.loginV1)
  router.Post("/v1/refresh", c.refreshV1)
}

// Perform a token refresh as mentioned in task
func (c authController) refreshV1(ctx *fiber.Ctx) error {
  access := ctx.Query("accessToken")
  refresh := ctx.Query("refreshToken")

  token, err := c.authService.RefreshV1(refresh, access, ctx.IP())
  if err != nil {
    return ctx.Status(400).JSON(fiber.Map{
      "status": "fail",
      "message": err.Error(),
    })
  }

  return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
    "status": "ok",
    "token": token,
  })
}

// Perform a tokens generation as mentioned in task
func (c authController) loginV1(ctx *fiber.Ctx) error {
  guid := ctx.Query("guid")

  token, err := c.authService.LoginV1(guid, ctx.IP())
  if err != nil {
    return ctx.Status(400).JSON(fiber.Map{
      "status": "fail",
      "message": err.Error(),
    })
  }

  return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
    "status": "ok",
    "token": token,
  })
}

func (c authController) login(ctx *fiber.Ctx) error {
  var payload *schemas.AuthLoginSchema
  
  if err := ctx.BodyParser(&payload); err != nil {
    return ctx.Status(422).JSON(fiber.Map{
      "status": "fail",
      "message": "Invalid body",
    })
  }

  token, err := c.authService.Login(payload, ctx.IP())
  if err != nil {
    return ctx.Status(400).JSON(fiber.Map{
      "status": "fail",
      "message": err.Error(),
    })
  }

  return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
    "status": "ok",
    "token": token,
  })
}

func (c authController) register(ctx *fiber.Ctx) error {
  var payload *schemas.AuthRegisterSchema
  
  if err := ctx.BodyParser(&payload); err != nil {
    return ctx.Status(422).JSON(fiber.Map{
      "status": "fail",
      "message": "Invalid body",
    })
  }

  user, err := c.authService.Register(payload)
  if err != nil {
    return ctx.Status(400).JSON(fiber.Map{
      "status": "fail",
      "message": err.Error(),
    })
  }

  return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
    "status": "ok",
    "user": user,
  })
}

func (c authController) refresh(ctx *fiber.Ctx) error {
  user := ctx.Locals("user").(*db.User)
  token := strings.TrimPrefix(ctx.Get("Authorization"), "Bearer ")

  tokens, err := c.authService.Refresh(user, token, ctx.IP())
  if err != nil {
    return ctx.Status(400).JSON(fiber.Map{
      "status": "fail",
      "message": err.Error(),
    })
  }

  return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
    "status": "ok",
    "token": tokens,
  })
}

func NewAuthController(authService services.AuthService) AuthController {
  return &authController{authService: authService}
}
