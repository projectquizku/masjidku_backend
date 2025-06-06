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

	userFollowMasjid "masjidku_backend/internals/features/masjids/user_follow_masjids/route"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MasjidPublicRoutes(r fiber.Router, db *gorm.DB) {
	// Ini endpoint yang boleh diakses publik tanpa login
	masjidRoutes.MasjidUserRoutes(r, db)
}

func MasjidUserRoutes(r fiber.Router, db *gorm.DB) {
	// Ini endpoint yang butuh login user biasa (dengan token)
	userFollowMasjid.UserFollowMasjidsRoutes(r, db)
	MasjidMore.MasjidMoreUserRoutes(r, db)
	LectureRoutes.UserLectureRoutes(r, db)
	EventRoutes.EventRoutesUser(r, db)
	FaqRoutes.FaqQuestionUserRoutes(r, db)
	LectureSessionRoutes.LectureSessionUserRoutes(r, db)
	LectureSessionsExamsRoutes.LectureSessionsExamsUserRoutes(r, db)
	LectureSessionsAssetRoutes.LectureSessionsAssetUserRoutes(r, db)
	LectureSessionsQuestionRoutes.LectureSessionsQuestionUserRoutes(r, db)
}

func MasjidAdminRoutes(r fiber.Router, db *gorm.DB) {
	// Ini endpoint khusus admin masjid
	masjidRoutes.MasjidAdminRoutes(r, db)
	MasjidAdmin.MasjidAdminRoutes(r, db)
	MasjidMore.MasjidMoreRoutes(r, db)
	LectureRoutes.LectureRoutes(r, db)
	EventRoutes.EventRoutes(r, db)
	NotificationRoutes.NotificationRoutes(r, db)
	FaqRoutes.FaqQuestionAdminRoutes(r, db)
	LectureSessionRoutes.LectureSessionAdminRoutes(r, db)
	LectureSessionsExamsRoutes.LectureSessionsExamsAdminRoutes(r, db)
	LectureSessionsAssetRoutes.LectureSessionsAssetAdminRoutes(r, db)
	LectureSessionsQuestionRoutes.LectureSessionsQuestionAdminRoutes(r, db)
}
