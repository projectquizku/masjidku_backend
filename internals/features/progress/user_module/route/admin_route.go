package submit_batch_route

import (
	"masjidku_backend/internals/features/progress/user_module/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubmitBatchAdminRoutes(api fiber.Router, db *gorm.DB) {
	ctrl := controller.NewSubmitBatchController(db)

	// Admin endpoints
	api.Get("/module-attempts", ctrl.GetAllModuleAttempts)
	api.Get("/answer-attempts", ctrl.GetAllAnswerAttempts)
}
