package route

import (
	unitController "masjidku_backend/internals/features/lessons/units/controller"
	"masjidku_backend/internals/constants"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UnitAdminRoutes(api fiber.Router, db *gorm.DB) {
	unitCtrl := unitController.NewUnitController(db)
	unitNewsCtrl := unitController.NewUnitNewsController(db)

	unitRoutes := api.Group("/units",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola unit"),
			constants.OwnerAndAbove,
		),
	)
	unitRoutes.Post("/", unitCtrl.CreateUnit)
	unitRoutes.Put("/:id", unitCtrl.UpdateUnit)
	unitRoutes.Delete("/:id", unitCtrl.DeleteUnit)

	unitNewsRoutes := api.Group("/units-news",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola berita unit"),
			constants.OwnerAndAbove,
		),
	)
	unitNewsRoutes.Post("/", unitNewsCtrl.Create)
	unitNewsRoutes.Put("/:id", unitNewsCtrl.Update)
	unitNewsRoutes.Delete("/:id", unitNewsCtrl.Delete)
}
