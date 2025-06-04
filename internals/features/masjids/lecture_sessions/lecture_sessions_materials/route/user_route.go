package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_materials/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// 👥 User Routes (Read-only)
func LectureSessionsAssetUserRoutes(user fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionsAssetController(db)

	user.Get("/lecture-sessions-assets", ctrl.GetAllLectureSessionsAssets)     // 📄 Lihat semua asset
	user.Get("/lecture-sessions-assets/:id", ctrl.GetLectureSessionsAssetByID) // 🔍 Detail asset

	ctrl2 := controller.NewLectureSessionsMaterialController(db)

	user.Get("/lecture-sessions-materials", ctrl2.GetAllLectureSessionsMaterials)     // 📄 Semua materi (read only)
	user.Get("/lecture-sessions-materials/:id", ctrl2.GetLectureSessionsMaterialByID) // 🔍 Detail materi
}
