package controllers

import (
  "github.com/gofiber/fiber/v2"
  "github.com/romakot321/go-jwt-api/internal/app/services"
  "github.com/romakot321/go-jwt-api/internal/app/middleware"
  "github.com/romakot321/go-jwt-api/internal/app/db"
)

type UserController interface {
  Register(router fiber.Router)
}

type userController struct {
  userService services.UserService
}

func (c userController) Register(router fiber.Router) {
  router.Post("/me", middleware.AuthenticateUserAccess, c.getMe)
}

func (c userController) getMe(ctx *fiber.Ctx) error {
  user := ctx.Locals("user").(*db.User)
  return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
    "status": "success",
    "user": user,
  })
}

func NewUserController(userService services.UserService) UserController {
  return &userController{userService: userService}
}
