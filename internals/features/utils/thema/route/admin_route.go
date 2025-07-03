package route

import (
	"masjidku_backend/internals/constants"
	themeController "masjidku_backend/internals/features/utils/thema/controller"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ThemeAdminRoutes(router fiber.Router, db *gorm.DB) {
	themeCtrl := themeController.NewThemeController(db)

	// 🔐 Route admin/teacher/owner
	adminRoutes := router.Group("/themes",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorNonUser("theme"),
			constants.TeacherAndAbove,
		),
	)

	adminRoutes.Get("/", themeCtrl.GetAllThemes)
	adminRoutes.Post("/", themeCtrl.CreateTheme)

}
