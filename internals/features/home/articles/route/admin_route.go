package route

import (
	"masjidku_backend/internals/features/home/articles/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ArticleAdminRoutes(api fiber.Router, db *gorm.DB) {
	articleCtrl := controller.NewArticleController(db)

	// === ADMIN ROUTES ===
	admin := api.Group("/articles")
	admin.Post("/", articleCtrl.CreateArticle)    // ➕ Buat artikel baru
	admin.Put("/:id", articleCtrl.UpdateArticle)  // 🔄 Perbarui artikel
	admin.Delete("/:id", articleCtrl.DeleteArticle) // 🗑️ Hapus artikel
}
