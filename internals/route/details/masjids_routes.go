package details

import (
	masjidRoutes "masjidku_backend/internals/features/masjids/masjids/route"
	rateLimiter "masjidku_backend/internals/middlewares"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MasjidRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// Group untuk admin: /api/a/...
	adminGroup := api.Group("/a")
	masjidRoutes.MasjidAdminRoutes(adminGroup, db)

	// Group untuk user/public: /api/u/...
	userGroup := api.Group("/u")
	masjidRoutes.MasjidUserRoutes(userGroup, db)

}
