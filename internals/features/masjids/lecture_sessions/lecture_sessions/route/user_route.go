package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ User Routes (read-only atau sesuai akses yang dibolehkan)
func LectureSessionUserRoutes(user fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionController(db)

	user.Get("/lecture-sessions", ctrl.GetAllLectureSessions)     // 📄 Lihat semua sesi
	user.Get("/lecture-sessions/:id", ctrl.GetLectureSessionByID) // 🔍 Detail sesi

	ctrl2 := controller.NewUserLectureSessionController(db)

	user.Post("/user-lecture-sessions", ctrl2.CreateUserLectureSession) // ✅ User mencatat kehadiran/progress
	user.Get("/user-lecture-sessions", ctrl2.GetAllUserLectureSessions) // 🔍 Lihat semua sesi yang diikuti
	user.Get("/user-lecture-sessions/:id", ctrl2.GetUserLectureSessionByID)
}
