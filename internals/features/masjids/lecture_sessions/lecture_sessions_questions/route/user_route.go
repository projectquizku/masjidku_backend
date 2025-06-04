package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_questions/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LectureSessionsQuestionUserRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewLectureSessionsQuestionController(db)
	api.Get("/", ctrl.GetAllLectureSessionsQuestions)
	api.Get("/:id", ctrl.GetLectureSessionsQuestionByID)

	ctrl2 := controller.NewLectureSessionsUserQuestionController(db)

	api.Post("/lecture-sessions-user-questions", ctrl2.CreateLectureSessionsUserQuestion)
	api.Get("/lecture-sessions-user-questions", ctrl2.GetAllLectureSessionsUserQuestions)
	api.Get("/lecture-sessions-user-questions/by-question/:question_id", ctrl2.GetByQuestionID)

}
