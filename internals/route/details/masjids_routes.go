package details

import (
	FaqRoutes "masjidku_backend/internals/features/home/faqs/route"
	NotificationRoutes "masjidku_backend/internals/features/home/notifications/route"
	LectureSessionsExamsRoutes "masjidku_backend/internals/features/masjids/lecture_sessions/exams/route"
	LectureSessionRoutes "masjidku_backend/internals/features/masjids/lecture_sessions/main/route"
	LectureSessionsAssetRoutes "masjidku_backend/internals/features/masjids/lecture_sessions/materials/route"
	LectureSessionsQuestionRoutes "masjidku_backend/internals/features/masjids/lecture_sessions/questions/route"
	EventRoutes "masjidku_backend/internals/features/masjids/lectures/events/route"
	LectureRoutes "masjidku_backend/internals/features/masjids/lectures/lectures/route"
	MasjidAdmin "masjidku_backend/internals/features/masjids/masjid_admins/route"
	masjidRoutes "masjidku_backend/internals/features/masjids/masjids/route"
	MasjidMore "masjidku_backend/internals/features/masjids/masjids_more/route"
	rateLimiter "masjidku_backend/internals/middlewares"
	authMiddleware "masjidku_backend/internals/middlewares/auth"
	MasjidkuMiddleware "masjidku_backend/internals/middlewares/features"

	userFollowMasjid "masjidku_backend/internals/features/masjids/user_follow_masjids/route"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MasjidRoutes(app *fiber.App, db *gorm.DB) {
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// Group untuk admin: /api/a/...
	adminGroup := api.Group("/a",
		MasjidkuMiddleware.IsMasjidAdmin(), // 💥 hanya grup ini yang cek masjid_id
	)
	masjidRoutes.MasjidAdminRoutes(adminGroup, db)
	MasjidAdmin.MasjidAdminRoutes(adminGroup, db)
	MasjidMore.MasjidMoreRoutes(adminGroup, db)
	LectureRoutes.LectureRoutes(adminGroup, db)
	EventRoutes.EventRoutes(adminGroup, db)
	NotificationRoutes.NotificationRoutes(adminGroup, db)
	FaqRoutes.FaqQuestionAdminRoutes(adminGroup, db)
	LectureSessionRoutes.LectureSessionAdminRoutes(adminGroup, db)
	LectureSessionsExamsRoutes.LectureSessionsExamsAdminRoutes(adminGroup, db)
	LectureSessionsAssetRoutes.LectureSessionsAssetAdminRoutes(adminGroup, db)
	LectureSessionsQuestionRoutes.LectureSessionsQuestionAdminRoutes(adminGroup, db)

	// Group untuk user/public: /api/u/...
	userGroup := api.Group("/u")
	masjidRoutes.MasjidUserRoutes(userGroup, db)
	userFollowMasjid.UserFollowMasjidsRoutes(userGroup, db)
	MasjidMore.MasjidMoreUserRoutes(userGroup, db)
	LectureRoutes.UserLectureRoutes(userGroup, db)
	EventRoutes.EventRoutesUser(userGroup, db)
	FaqRoutes.FaqQuestionUserRoutes(userGroup, db)
	LectureSessionRoutes.LectureSessionUserRoutes(userGroup, db)
	LectureSessionsExamsRoutes.LectureSessionsExamsUserRoutes(userGroup, db)
	LectureSessionsAssetRoutes.LectureSessionsAssetUserRoutes(userGroup, db)
	LectureSessionsQuestionRoutes.LectureSessionsQuestionUserRoutes(userGroup, db)

}
