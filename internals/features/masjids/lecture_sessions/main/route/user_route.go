package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ User Routes (read-only atau sesuai akses yang dibolehkan)
func LectureSessionUserRoutes(user fiber.Router, db *gorm.DB) {
	lectureSessionCtrl := controller.NewLectureSessionController(db)
	userLectureSessionCtrl := controller.NewUserLectureSessionController(db)

	// 📚 Group: /lecture-sessions
	session := user.Group("/lecture-sessions")
	session.Get("/", lectureSessionCtrl.GetAllLectureSessions)    // 📄 Lihat semua sesi
	session.Get("/:id", lectureSessionCtrl.GetLectureSessionByID) // 🔍 Detail sesi

	// 👥 Group: /user-lecture-sessions
	userSession := user.Group("/user-lecture-sessions")
	userSession.Post("/", userLectureSessionCtrl.CreateUserLectureSession)    // ✅ Catat kehadiran / progress
	userSession.Get("/", userLectureSessionCtrl.GetAllUserLectureSessions)    // 🔍 Lihat semua sesi yang diikuti
	userSession.Get("/:id", userLectureSessionCtrl.GetUserLectureSessionByID) // 🔍 Detail kehadiran
}
