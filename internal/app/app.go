package app

import (
  "log"

  "github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/cors"
  "github.com/romakot321/go-jwt-api/internal/app/controllers"
  "github.com/romakot321/go-jwt-api/internal/app/repositories"
  "github.com/romakot321/go-jwt-api/internal/app/services"
)

func Run() {
  authRepository := repositories.NewAuthRepository("REPLACEME")
  userRepository := repositories.NewUserRepository()
  tokenRepository := repositories.NewTokenRepository()
  mailRepository := repositories.NewMailRepository()

  authService := services.NewAuthService(authRepository, userRepository, tokenRepository, mailRepository)
  userService := services.NewUserService(userRepository)

  authController := controllers.NewAuthController(authService)
  userController := controllers.NewUserController(userService)

  app := fiber.New()
  router := fiber.New()

  app.Mount("/api", router)
  app.Use(cors.New(cors.Config{
    AllowOrigins: "http://localhost:8000",
    AllowHeaders: "*",
    AllowMethods: "*",
    AllowCredentials: true,
  }))

  router.Route("/auth", authController.Register)
  router.Route("/user", userController.Register)

  router.All("*", func(c *fiber.Ctx) error {
    return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "not found"})
  })

  log.Fatal(app.Listen(":8000"))
}
