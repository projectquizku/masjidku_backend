package route

import (
	"masjidku_backend/internals/features/masjids/masjids/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MasjidAdminRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewMasjidController(db)

	// Hanya admin yang bisa melakukan ini:
	api.Post("/masjids", ctrl.CreateMasjid)
	api.Put("/masjids/:id", ctrl.UpdateMasjid)
	api.Delete("/masjids/:id", ctrl.DeleteMasjid)

	ctrl2 := controller.NewMasjidProfileController(db)

	// Hanya admin yang bisa melakukan ini:
	api.Post("/masjid-profiles", ctrl2.CreateMasjidProfile)
	api.Put("/masjid-profiles/:id", ctrl2.UpdateMasjidProfile)
	api.Delete("/masjid-profiles/:id", ctrl2.DeleteMasjidProfile)
}
