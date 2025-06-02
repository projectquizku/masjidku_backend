package route

import (
	difficultyController "masjidku_backend/internals/features/lessons/difficulty/controller"
	"masjidku_backend/internals/constants"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func DifficultyAdminRoutes(api fiber.Router, db *gorm.DB) {
	difficultyCtrl := difficultyController.NewDifficultyController(db)
	difficultyNewsCtrl := difficultyController.NewDifficultyNewsController(db)

	// 🔐 Role admin, teacher, owner
	protectedDifficultyRoutes := api.Group("/difficulties",
		authMiddleware.OnlyRoles(
			constants.RoleErrorTeacher("mengelola difficulties"),
			constants.AdminOnly...,
		),
	)
	protectedDifficultyRoutes.Post("/", difficultyCtrl.CreateDifficulty)
	protectedDifficultyRoutes.Put("/:id", difficultyCtrl.UpdateDifficulty)
	protectedDifficultyRoutes.Delete("/:id", difficultyCtrl.DeleteDifficulty)

	protectedDifficultyNewsRoutes := api.Group("/difficulties-news",
		authMiddleware.OnlyRoles(
			constants.RoleErrorTeacher("mengelola difficulty news"),
			constants.OwnerAndAbove...,
		),
	)
	protectedDifficultyNewsRoutes.Post("/", difficultyNewsCtrl.CreateNews)
	protectedDifficultyNewsRoutes.Put("/:id", difficultyNewsCtrl.UpdateNews)
	protectedDifficultyNewsRoutes.Delete("/:id", difficultyNewsCtrl.DeleteNews)
}
