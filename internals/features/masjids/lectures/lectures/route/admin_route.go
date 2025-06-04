package route

import (
	"masjidku_backend/internals/features/masjids/lectures/lectures/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LectureRoutes(router fiber.Router, db *gorm.DB) {
	ctrl := controller.LectureController{DB: db}

	router.Post("/lectures", ctrl.CreateLecture)
	router.Get("/lectures", ctrl.GetLecturesByMasjid)

	ctrl2 := controller.UserLectureController{DB: db}

	router.Post("/user-lectures", ctrl2.CreateUserLecture)
	router.Get("/user-lectures", ctrl2.GetUsersByLecture)
}
