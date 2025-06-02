package route

import (
	examController "masjidku_backend/internals/features/quizzes/exams/controller"
	"masjidku_backend/internals/constants"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ExamAdminRoutes(api fiber.Router, db *gorm.DB) {
	examCtrl := examController.NewExamController(db)

	examRoutes := api.Group("/exams",
		authMiddleware.OnlyRolesSlice(
			constants.RoleErrorTeacher("mengelola ujian"),
			constants.TeacherAndAbove,
		),
	)
	examRoutes.Post("/", examCtrl.CreateExam)
	examRoutes.Put("/:id", examCtrl.UpdateExam)
	examRoutes.Delete("/:id", examCtrl.DeleteExam)
}
