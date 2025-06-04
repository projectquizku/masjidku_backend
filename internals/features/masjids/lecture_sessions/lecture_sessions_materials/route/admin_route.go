package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_materials/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// 🔐 Admin Routes (CRUD)
func LectureSessionsAssetAdminRoutes(admin fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionsAssetController(db)

	admin.Post("/lecture-sessions-assets", ctrl.CreateLectureSessionsAsset)       // ➕ Tambah asset
	admin.Get("/lecture-sessions-assets", ctrl.GetAllLectureSessionsAssets)       // 📄 Lihat semua asset
	admin.Get("/lecture-sessions-assets/:id", ctrl.GetLectureSessionsAssetByID)   // 🔍 Detail asset
	admin.Delete("/lecture-sessions-assets/:id", ctrl.DeleteLectureSessionsAsset) // ❌ Hapus asset

	ctrl2 := controller.NewLectureSessionsMaterialController(db)

	admin.Post("/lecture-sessions-materials", ctrl2.CreateLectureSessionsMaterial)       // ➕ Tambah materi
	admin.Get("/lecture-sessions-materials", ctrl2.GetAllLectureSessionsMaterials)       // 📄 Semua materi
	admin.Get("/lecture-sessions-materials/:id", ctrl2.GetLectureSessionsMaterialByID)   // 🔍 Detail
	admin.Delete("/lecture-sessions-materials/:id", ctrl2.DeleteLectureSessionsMaterial) // ❌ Hapus
}
