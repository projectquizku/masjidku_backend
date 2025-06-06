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
	// authMiddleware "masjidku_backend/internals/middlewares/auth"
	MasjidkuMiddleware "masjidku_backend/internals/middlewares/features"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func HomeRoutes(app *fiber.App, db *gorm.DB) {
	// Grup umum (tanpa admin masjid check)
	api := app.Group("/api",
		// authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// Grup ADMIN (dengan validasi masjid_admin_ids via JWT)
	adminDKMGroup := api.Group("/a",
		MasjidkuMiddleware.IsMasjidAdmin(), // 💥 hanya grup ini yang cek masjid_id
	)
	// Semua /api/a/... route hanya bisa diakses jika user adalah admin masjid
	FaqRoutes.FaqQuestionAdminRoutes(adminDKMGroup, db)
	AdviceRoutes.AdviceAdminRoutes(adminDKMGroup, db)
	ArticleRoutes.ArticleAdminRoutes(adminDKMGroup, db)
	PostRoutes.PostAdminRoutes(adminDKMGroup, db)
	QuestionnaireRoutes.QuestionnaireQuestionAdminRoutes(adminDKMGroup, db)

	adminGroup := api.Group("/a")
	QouteRoutes.QuoteAdminRoutes(adminGroup, db)


	// Grup USER biasa (cukup login + rate limiter)
	userGroup := api.Group("/u")
	NotificationRoutes.NotificationUserRoutes(userGroup, db)
	FaqRoutes.FaqQuestionUserRoutes(userGroup, db)
	AdviceRoutes.AdviceUserRoutes(userGroup, db)
	ArticleRoutes.ArticleUserRoutes(userGroup, db)
	QouteRoutes.QuoteUserRoutes(userGroup, db)
	PostRoutes.PostUserRoutes(userGroup, db)
	QuestionnaireRoutes.QuestionnaireQuestionUserRoutes(userGroup, db)
}
