package route

import (
	categoryController "masjidku_backend/internals/features/lessons/categories/controller"
	"masjidku_backend/internals/constants"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CategoryAdminRoutes(api fiber.Router, db *gorm.DB) {
	categoryCtrl := categoryController.NewCategoryController(db)
	categoryNewsCtrl := categoryController.NewCategoryNewsController(db)

	categoryRoutes := api.Group("/categories",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola kategori"),
			constants.OwnerAndAbove,
		),
	)
	categoryRoutes.Post("/", categoryCtrl.CreateCategory)
	categoryRoutes.Put("/:id", categoryCtrl.UpdateCategory)
	categoryRoutes.Delete("/:id", categoryCtrl.DeleteCategory)

	categoryNewsRoutes := api.Group("/categories-news",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola berita kategori"),
			constants.OwnerAndAbove,
		),
	)
	categoryNewsRoutes.Post("/", categoryNewsCtrl.Create)
	categoryNewsRoutes.Put("/:id", categoryNewsCtrl.Update)
	categoryNewsRoutes.Delete("/:id", categoryNewsCtrl.Delete)
}
