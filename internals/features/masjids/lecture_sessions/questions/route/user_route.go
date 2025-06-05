package route

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/questions/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func LectureSessionsQuestionUserRoutes(user fiber.Router, db *gorm.DB) {
	questionCtrl := controller.NewLectureSessionsQuestionController(db)
	userQuestionCtrl := controller.NewLectureSessionsUserQuestionController(db)

	// 📝 Group: /lecture-sessions-questions (read-only)
	questions := user.Group("/lecture-sessions-questions")
	questions.Get("/", questionCtrl.GetAllLectureSessionsQuestions)
	// questions.Get("/:id", questionCtrl.GetLectureSessionsQuestionByID) // (opsional)

	// 👤 Group: /lecture-sessions-user-questions
	userQuestions := user.Group("/lecture-sessions-user-questions")
	userQuestions.Post("/", userQuestionCtrl.CreateLectureSessionsUserQuestion)
	// userQuestions.Get("/", userQuestionCtrl.GetAllUserLectureSessionsQuestions)
	userQuestions.Get("/by-question/:question_id", userQuestionCtrl.GetByQuestionID)
	// userQuestions.Get("/:id", userQuestionCtrl.GetLectureSessionsUserQuestionByID) // (opsional)
}
