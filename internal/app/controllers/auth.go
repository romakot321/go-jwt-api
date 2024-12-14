package controllers

import (
  "github.com/gofiber/fiber/v2"
  "github.com/romakot321/go-jwt-api/internal/app/schemas"
  "github.com/romakot321/go-jwt-api/internal/app/services"
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
}

func (c authController) login(ctx *fiber.Ctx) error {
  var payload *schemas.AuthLoginSchema
  
  if err := ctx.BodyParser(&payload); err != nil {
    return ctx.Status(422).JSON(fiber.Map{
      "status": "fail",
      "message": "Invalid body",
    })
  }

  token, err := c.authService.Login(payload)
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

func NewAuthController(authService services.AuthService) AuthController {
  return &authController{authService: authService}
}
