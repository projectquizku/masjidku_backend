package submit_batch_route

import (
	"masjidku_backend/internals/features/progress/user_module/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SubmitBatchUserRoutes(api fiber.Router, db *gorm.DB) {
	submitBatchCtrl := controller.NewSubmitBatchController(db)

	// Endpoint utama untuk kirim seluruh modul attempt secara batch
	api.Post("/submit-batch", submitBatchCtrl.Submit)
	api.All("/get-batch/:batch_id", submitBatchCtrl.GetBatch)
}
