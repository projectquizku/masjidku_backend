package details

import (
	categoriesRoute "masjidku_backend/internals/features/lessons/categories/route"
	difficultiesRoute "masjidku_backend/internals/features/lessons/difficulty/route"
	subcategoriesRoute "masjidku_backend/internals/features/lessons/subcategories/route"
	themesOrLevelsRoute "masjidku_backend/internals/features/lessons/themes_or_levels/route"
	unitsRoute "masjidku_backend/internals/features/lessons/units/route"

	rateLimiter "masjidku_backend/internals/middlewares"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LessonRoutes(app *fiber.App, db *gorm.DB) {
	// 🔐 Semua lesson route pakai auth & rate limiter
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 👤 Prefix untuk user biasa
	userGroup := api.Group("/u")
	difficultiesRoute.DifficultyUserRoutes(userGroup, db)
	categoriesRoute.CategoryUserRoutes(userGroup, db)
	subcategoriesRoute.SubcategoryUserRoutes(userGroup, db)
	themesOrLevelsRoute.ThemesOrLevelsUserRoutes(userGroup, db)
	unitsRoute.UnitUserRoutes(userGroup, db)

	// 🔐 Prefix untuk admin
	adminGroup := api.Group("/a")
	difficultiesRoute.DifficultyAdminRoutes(adminGroup, db)
	categoriesRoute.CategoryAdminRoutes(adminGroup, db)
	subcategoriesRoute.SubcategoryAdminRoutes(adminGroup, db)
	themesOrLevelsRoute.ThemesOrLevelsAdminRoutes(adminGroup, db)
	unitsRoute.UnitAdminRoutes(adminGroup, db)
}
