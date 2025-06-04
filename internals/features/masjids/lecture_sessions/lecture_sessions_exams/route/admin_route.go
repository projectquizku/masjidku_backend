package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_exams/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LectureSessionsExamsAdminRoutes(admin fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionsExamController(db)

	admin.Post("/lecture-sessions-exams", ctrl.CreateLectureSessionsExam)       // ➕ Buat ujian sesi kajian
	admin.Get("/lecture-sessions-exams", ctrl.GetAllLectureSessionsExams)       // 📄 Lihat semua ujian
	admin.Get("/lecture-sessions-exams/:id", ctrl.GetLectureSessionsExamByID)   // 🔍 Detail ujian
	admin.Put("/lecture-sessions-exams/:id", ctrl.UpdateLectureSessionsExam)    // ✏️ Edit ujian
	admin.Delete("/lecture-sessions-exams/:id", ctrl.DeleteLectureSessionsExam) // ❌ Hapus ujian

	ctrl2 := controller.NewUserLectureSessionsExamController(db)

	admin.Get("/user-lecture-sessions-exams", ctrl2.GetAllUserLectureSessionsExams) // 📄 Lihat semua
	admin.Get("/user-lecture-sessions-exams/:id", ctrl2.GetUserLectureSessionsExamByID)
}
