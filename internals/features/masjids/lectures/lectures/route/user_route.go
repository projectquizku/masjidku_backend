package route

import (
	"masjidku_backend/internals/features/masjids/lectures/lectures/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func UserLectureRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewUserLectureController(db)

	lecture := api.Group("/user-lectures")
	lecture.Post("/", ctrl.CreateUserLecture)
	lecture.Post("/by-lecture", ctrl.GetUsersByLecture)// ✅ opsional tambahan jika ingin ambil semua kajian yang diikuti user
}
