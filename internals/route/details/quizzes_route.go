package details

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	rateLimiter "masjidku_backend/internals/middlewares"
	authMiddleware "masjidku_backend/internals/middlewares/auth"

	evaluationRoute "masjidku_backend/internals/features/quizzes/evaluations/route"
	examsRoute "masjidku_backend/internals/features/quizzes/exams/route"
	questionsRoute "masjidku_backend/internals/features/quizzes/questions/route"
	quizzesRoute "masjidku_backend/internals/features/quizzes/quizzes/route"
	readingsRoute "masjidku_backend/internals/features/quizzes/readings/route"
)

func QuizzesRoute(app *fiber.App, db *gorm.DB) {
	// Bungkus dengan Auth dan RateLimiter
	api := app.Group("/api",
		authMiddleware.AuthMiddleware(db),
		rateLimiter.GlobalRateLimiter(),
	)

	// 👤 Prefix user: /api/u
	userGroup := api.Group("/u")
	quizzesRoute.QuizzesUserRoutes(userGroup, db)
	evaluationRoute.EvaluationUserRoutes(userGroup, db)
	examsRoute.ExamUserRoutes(userGroup, db)
	readingsRoute.ReadingUserRoutes(userGroup, db)
	questionsRoute.QuestionUserRoutes(userGroup, db)

	// 🔐 Prefix admin: /api/a
	adminGroup := api.Group("/a")
	quizzesRoute.QuizzesAdminRoutes(adminGroup, db)
	evaluationRoute.EvaluationAdminRoutes(adminGroup, db)
	examsRoute.ExamAdminRoutes(adminGroup, db)
	readingsRoute.ReadingAdminRoutes(adminGroup, db)
	questionsRoute.QuestionAdminRoutes(adminGroup, db)
}
