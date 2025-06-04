package route

import (
	"masjidku_backend/internals/features/home/faqs/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FaqQuestionAdminRoutes(admin fiber.Router, db *gorm.DB) {
	ctrl := controller.NewFaqQuestionController(db)

	admin.Get("/faq-questions", ctrl.GetAllFaqQuestions)       // ✅ Lihat semua pertanyaan
	admin.Get("/faq-questions/:id", ctrl.GetFaqQuestionByID)   // ✅ Detail pertanyaan
	admin.Put("/faq-questions/:id", ctrl.UpdateFaqQuestion)    // ✅ Tandai sebagai dijawab / edit
	admin.Delete("/faq-questions/:id", ctrl.DeleteFaqQuestion) // ✅ Hapus pertanyaan

	ctrl2 := controller.NewFaqAnswerController(db)

	admin.Post("/faq-answers", ctrl2.CreateFaqAnswer)       // ✅ Admin menjawab pertanyaan
	admin.Put("/faq-answers/:id", ctrl2.UpdateFaqAnswer)    // ✅ Edit jawaban
	admin.Delete("/faq-answers/:id", ctrl2.DeleteFaqAnswer) // ✅ Hapus jawaban
	admin.Get("/faq-answers/:id", ctrl2.GetFaqAnswerByID)   // ✅ Lihat detail jawaban
}
