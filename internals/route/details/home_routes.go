package details

import (
	AdviceRoutes "masjidku_backend/internals/features/home/advices/route"
	ArticleRoutes "masjidku_backend/internals/features/home/articles/route"
	FaqRoutes "masjidku_backend/internals/features/home/faqs/route"
	NotificationRoutes "masjidku_backend/internals/features/home/notifications/route"
	PostRoutes "masjidku_backend/internals/features/home/posts/route"
	QouteRoutes "masjidku_backend/internals/features/home/qoutes/route"
	QuestionnaireRoutes "masjidku_backend/internals/features/home/questionnaires/route"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ Untuk route publik tanpa token
// Contoh akses: /api/public/quotes
func HomePublicRoutes(api fiber.Router, db *gorm.DB) {
	QouteRoutes.QuoteUserRoutes(api, db)
	AdviceRoutes.AdviceUserRoutes(api, db)
	FaqRoutes.FaqQuestionUserRoutes(api, db)
	ArticleRoutes.ArticleUserRoutes(api, db)
	PostRoutes.PostUserRoutes(api, db)
	QuestionnaireRoutes.QuestionnaireQuestionUserRoutes(api, db)
}

// ✅ Untuk route user login (dengan token)
// Contoh akses: /api/u/notifications
func HomePrivateRoutes(api fiber.Router, db *gorm.DB) {
	private := api.Group("/u")
	NotificationRoutes.NotificationUserRoutes(private, db)
}

// ✅ Untuk route admin masjid (token + admin)
// Contoh akses: /api/a/quotes
func HomeAdminRoutes(api fiber.Router, db *gorm.DB) {
	FaqRoutes.FaqQuestionAdminRoutes(api, db)
	AdviceRoutes.AdviceAdminRoutes(api, db)
	ArticleRoutes.ArticleAdminRoutes(api, db)
	PostRoutes.PostAdminRoutes(api, db)
	QuestionnaireRoutes.QuestionnaireQuestionAdminRoutes(api, db)
	QouteRoutes.QuoteAdminRoutes(api, db)
}
