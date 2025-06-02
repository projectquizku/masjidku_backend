package route

import (
	readingController "masjidku_backend/internals/features/quizzes/readings/controller"
	"masjidku_backend/internals/constants"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReadingAdminRoutes(api fiber.Router, db *gorm.DB) {
	readingCtrl := readingController.NewReadingController(db)

	readingRoutes := api.Group("/readings",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola reading"),
			constants.TeacherAndAbove,
		),
	)
	readingRoutes.Post("/", readingCtrl.CreateReading)
	readingRoutes.Put("/:id", readingCtrl.UpdateReading)
	readingRoutes.Delete("/:id", readingCtrl.DeleteReading)
}
