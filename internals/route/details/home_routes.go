package details

import (
	AdviceRoutes "masjidku_backend/internals/features/home/advices/route"
	ArticleRoutes "masjidku_backend/internals/features/home/articles/route"
	FaqRoutes "masjidku_backend/internals/features/home/faqs/route"
	NotificationRoutes "masjidku_backend/internals/features/home/notifications/route"
	PostRoutes "masjidku_backend/internals/features/home/posts/route"
	QouteRoutes "masjidku_backend/internals/features/home/qoutes/route"
	QuestionnaireRoutes "masjidku_backend/internals/features/home/questionnaires/route"
	rateLimiter "masjidku_backend/internals/middlewares"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HomeRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// Group untuk admin: /api/a/...
	adminGroup := api.Group("/a")

	NotificationRoutes.NotificationRoutes(adminGroup, db)
	FaqRoutes.FaqQuestionAdminRoutes(adminGroup, db)
	AdviceRoutes.AdviceAdminRoutes(adminGroup, db)
	ArticleRoutes.ArticleAdminRoutes(adminGroup, db)
	QouteRoutes.QuoteAdminRoutes(adminGroup, db)
	PostRoutes.PostAdminRoutes(adminGroup, db)
	QuestionnaireRoutes.QuestionnaireQuestionAdminRoutes(adminGroup, db)
	// Group untuk user/public: /api/u/...

	adminGroup = api.Group("/u")
	NotificationRoutes.NotificationUserRoutes(adminGroup, db)
	FaqRoutes.FaqQuestionUserRoutes(adminGroup, db)
	AdviceRoutes.AdviceUserRoutes(adminGroup, db)
	ArticleRoutes.ArticleUserRoutes(adminGroup, db)
	QouteRoutes.QuoteUserRoutes(adminGroup, db)
	PostRoutes.PostUserRoutes(adminGroup, db)
	QuestionnaireRoutes.QuestionnaireQuestionUserRoutes(adminGroup, db)
}
