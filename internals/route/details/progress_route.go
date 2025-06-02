package details

import (
	levelRankRoute "masjidku_backend/internals/features/progress/level_rank/route"
	userPointRoutes "masjidku_backend/internals/features/progress/points/route"
	userProgressRoutes "masjidku_backend/internals/features/progress/progress/route"
	submitBatchRoutes "masjidku_backend/internals/features/progress/user_module/route"
	rateLimiter "masjidku_backend/internals/middlewares"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ProgressRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	adminGroup := api.Group("/a")
	levelRankRoute.LevelRequirementAdminRoute(adminGroup, db)
	submitBatchRoutes.SubmitBatchAdminRoutes(adminGroup, db)

	userGroup := api.Group("/u")
	levelRankRoute.LevelRequirementUserRoute(userGroup, db)

	// ✅ Ini diperbaiki
	userPointRoutes.UserPointRoutes(userGroup, db)
	userProgressRoutes.UserProgressRoutes(userGroup, db)
	submitBatchRoutes.SubmitBatchUserRoutes(userGroup, db)
}
