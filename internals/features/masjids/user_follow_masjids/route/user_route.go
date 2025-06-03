package route

import (
	"masjidku_backend/internals/features/masjids/user_follow_masjids/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserFollowMasjidsRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewUserFollowMasjidController(db)

	// 🟢 User follow dan unfollow masjid
	api.Post("/user-follow-masjids/follow", ctrl.FollowMasjid)
	api.Delete("/user-follow-masjids/unfollow", ctrl.UnfollowMasjid)

	// 🟢 Get daftar masjid yang di-follow oleh user
	api.Get("/user-follow-masjids/followed", ctrl.GetFollowedMasjidsByUser) // Ambil dari body atau query
}
