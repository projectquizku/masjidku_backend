package details

import (
	masjidRoutes "masjidku_backend/internals/features/masjids/masjids/route"
	rateLimiter "masjidku_backend/internals/middlewares"
	authMiddleware "masjidku_backend/internals/middlewares/auth"
	MasjidAdmin "masjidku_backend/internals/features/masjids/masjid_admins/route"
	MasjidMore "masjidku_backend/internals/features/masjids/masjids_more/route"
	LectureRoutes "masjidku_backend/internals/features/masjids/lectures/lectures/route"

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
	adminGroup := api.Group("/a")
	masjidRoutes.MasjidAdminRoutes(adminGroup, db)
	MasjidAdmin.MasjidAdminRoutes(adminGroup, db)
	MasjidMore.MasjidMoreRoutes(adminGroup, db)
	LectureRoutes.LectureRoutes(adminGroup, db)

	// Group untuk user/public: /api/u/...
	userGroup := api.Group("/u")
	masjidRoutes.MasjidUserRoutes(userGroup, db)
	userFollowMasjid.UserFollowMasjidsRoutes(userGroup, db)

}
