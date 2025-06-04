package route

import (
	"masjidku_backend/internals/features/home/faqs/controller"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FaqQuestionUserRoutes(user fiber.Router, db *gorm.DB) {
	ctrl := controller.NewFaqQuestionController(db)

	user.Post("/faq-questions", ctrl.CreateFaqQuestion)     // ✅ Kirim pertanyaan
	user.Get("/faq-questions", ctrl.GetAllFaqQuestions)     // ✅ Lihat semua (bisa difilter nanti per user ID)
	user.Get("/faq-questions/:id", ctrl.GetFaqQuestionByID) // ✅ Detail pertanyaan

	ctrl2 := controller.NewFaqAnswerController(db)

	user.Get("/faq-answers/:id", ctrl2.GetFaqAnswerByID) // ✅ Lihat jawaban untuk pertanyaan tertentu

}
