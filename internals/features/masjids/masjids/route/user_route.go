package route

import (
	"masjidku_backend/internals/features/masjids/masjids/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MasjidUserRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewMasjidController(db)

	// User & publik bisa akses ini
	api.Get("/masjids", ctrl.GetAllMasjids)
	api.Get("/masjids/:slug", ctrl.GetMasjidBySlug)

	ctrl2 := controller.NewMasjidProfileController(db)

	// User & publik bisa akses ini
	api.Get("/masjid-profiles", ctrl2.GetProfileByMasjidID)
}
