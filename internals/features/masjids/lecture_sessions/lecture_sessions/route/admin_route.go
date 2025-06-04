package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ Admin Routes
func LectureSessionAdminRoutes(admin fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionController(db)

	admin.Post("/lecture-sessions", ctrl.CreateLectureSession)       // ➕ Buat sesi baru
	admin.Get("/lecture-sessions", ctrl.GetAllLectureSessions)       // 📄 Lihat semua sesi
	admin.Get("/lecture-sessions/:id", ctrl.GetLectureSessionByID)   // 🔍 Detail sesi
	admin.Put("/lecture-sessions/:id", ctrl.UpdateLectureSession)    // ✏️ Edit sesi
	admin.Delete("/lecture-sessions/:id", ctrl.DeleteLectureSession) // ❌ Hapus sesi

	ctrl2 := controller.NewUserLectureSessionController(db)

	admin.Get("/user-lecture-sessions", ctrl2.GetAllUserLectureSessions)
	admin.Get("/user-lecture-sessions/:id", ctrl2.GetUserLectureSessionByID)
	admin.Put("/user-lecture-sessions/:id", ctrl2.UpdateUserLectureSession)
	admin.Delete("/user-lecture-sessions/:id", ctrl2.DeleteUserLectureSession)
}
