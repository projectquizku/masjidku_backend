package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_exams/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LectureSessionsExamsUserRoutes(user fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionsExamController(db)

	user.Get("/lecture-sessions-exams", ctrl.GetAllLectureSessionsExams)     // 📄 Lihat semua ujian
	user.Get("/lecture-sessions-exams/:id", ctrl.GetLectureSessionsExamByID) // 🔍 Lihat detail ujian

	ctrl2 := controller.NewUserLectureSessionsExamController(db)

	user.Post("/user-lecture-sessions-exams", ctrl2.CreateUserLectureSessionsExam) // ➕ Kirim progress
	user.Get("/user-lecture-sessions-exams", ctrl2.GetAllUserLectureSessionsExams) // 📄 Lihat semua (nanti bisa difilter by user_id)
	user.Get("/user-lecture-sessions-exams/:id", ctrl2.GetUserLectureSessionsExamByID)
}
