package route

import (
	controller "masjidku_backend/internals/features/users/auth/controller"
	rateLimiter "masjidku_backend/internals/middlewares"
	authMw "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func AuthRoutes(app *fiber.App, db *gorm.DB) {
	// ✅ Inisialisasi semua controller
	authController := controller.NewAuthController(db)

	// ✅ Pasang Global Rate Limiter (hanya limit, tidak cek token)
	app.Use(rateLimiter.GlobalRateLimiter())

	// ✅ PUBLIC routes (TIDAK pakai JWT AuthMiddleware)
	publicAuth := app.Group("/auth")

	publicAuth.Post("/login", rateLimiter.LoginRateLimiter(), authController.Login)
	publicAuth.Options("/login", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	publicAuth.Post("/register", rateLimiter.RegisterRateLimiter(), authController.Register)
	publicAuth.Options("/register", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	publicAuth.Post("/forgot-password/check", rateLimiter.ForgotPasswordRateLimiter(), authController.CheckSecurityAnswer)
	publicAuth.Options("/forgot-password/check", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	publicAuth.Post("/forgot-password/reset", authController.ResetPassword)
	publicAuth.Options("/forgot-password/reset", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	publicAuth.Post("/login-google", authController.LoginGoogle)
	publicAuth.Options("/login-google", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	publicAuth.Post("/refresh-token", authController.RefreshToken)
	publicAuth.Options("/refresh-token", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	// ✅ PROTECTED routes (HARUS pakai JWT AuthMiddleware)
	protectedAuth := app.Group("/api/auth", authMw.AuthMiddleware(db))

	protectedAuth.Post("/logout", authController.Logout)
	protectedAuth.Options("/logout", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	protectedAuth.Post("/change-password", authController.ChangePassword)
	protectedAuth.Options("/change-password", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })
}
