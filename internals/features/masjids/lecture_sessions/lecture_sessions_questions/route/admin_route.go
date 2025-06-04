package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_questions/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LectureSessionsQuestionAdminRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionsQuestionController(db)

	api.Post("/", ctrl.CreateLectureSessionsQuestion)      // ➕ Tambah soal
	api.Get("/", ctrl.GetAllLectureSessionsQuestions)      // 📄 Lihat semua soal
	api.Get("/:id", ctrl.GetLectureSessionsQuestionByID)   // 🔍 Soal tertentu
	api.Delete("/:id", ctrl.DeleteLectureSessionsQuestion) // ❌ Hapus soal

	ctrl2 := controller.NewLectureSessionsUserQuestionController(db)

	api.Delete("/lecture-sessions-user-questions/:id", ctrl2.DeleteByID)
}
