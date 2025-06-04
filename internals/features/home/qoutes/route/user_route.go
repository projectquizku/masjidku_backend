package route

import (
	"masjidku_backend/internals/features/home/qoutes/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func QuoteUserRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewQuoteController(db)

	// === USER ROUTES ===
	user := api.Group("/quotes")
	user.Get("/", ctrl.GetAllQuotes)    // 📄 Lihat semua quote
	user.Get("/:id", ctrl.GetQuoteByID) // 🔍 Detail quote
}
