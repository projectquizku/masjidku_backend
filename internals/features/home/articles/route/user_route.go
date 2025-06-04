package route

import (
	"masjidku_backend/internals/features/home/articles/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ArticleUserRoutes(api fiber.Router, db *gorm.DB) {
	articleCtrl := controller.NewArticleController(db)

	// === USER ROUTES ===
	user := api.Group("/articles")
	user.Get("/", articleCtrl.GetAllArticles)     // 📄 Lihat semua artikel
	user.Get("/:id", articleCtrl.GetArticleByID)  // 🔍 Lihat detail artikel
}
