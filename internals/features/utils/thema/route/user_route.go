package route

import (
	themeController "masjidku_backend/internals/features/utils/thema/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ThemePublicRoutes(router fiber.Router, db *gorm.DB) {
	themeCtrl := themeController.NewThemeController(db)

	publicRoutes := router.Group("/themes")
	publicRoutes.Get("/", themeCtrl.GetAllThemes) // bebas tanpa login

	userThemeCtrl := themeController.NewUserThemeController(db)

	userRoutes := router.Group("/user-themes")
	userRoutes.Get("/", userThemeCtrl.GetUserThemes)            // semua tema user
	userRoutes.Get("/selected", userThemeCtrl.GetSelectedTheme) // tema aktif user
	userRoutes.Post("/unlock", userThemeCtrl.UnlockTheme)       // unlock tema
	userRoutes.Post("/select", userThemeCtrl.SelectTheme)
	userRoutes.Get("/with-fallback", userThemeCtrl.GetUserThemesWithFallback)
}
