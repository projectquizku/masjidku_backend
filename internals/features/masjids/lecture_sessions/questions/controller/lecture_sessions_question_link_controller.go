package controller

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/questions/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/questions/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ✅ Inisialisasi validator
var validate = validator.New()

type LectureSessionsQuestionLinkController struct {
	DB *gorm.DB
}

func NewLectureSessionsQuestionLinkController(db *gorm.DB) *LectureSessionsQuestionLinkController {
	return &LectureSessionsQuestionLinkController{DB: db}
}

// =============================
// ➕ Create Question Link (Quiz/Exam)
// =============================
func (ctrl *LectureSessionsQuestionLinkController) CreateLink(c *fiber.Ctx) error {
	var body dto.CreateLectureSessionsQuestionLinkRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Validasi pakai validator (pastikan struct validator kamu inisialisasi ya)
	if err := validate.Struct(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validasi logika minimal salah satu: exam_id atau quiz_id harus ada
	if body.ExamID == nil && body.QuizID == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Either exam_id or quiz_id must be provided")
	}

	link := model.LectureSessionsQuestionLinkModel{
		QuestionID:    body.QuestionID,
		ExamID:        body.ExamID,
		QuizID:        body.QuizID,
		QuestionOrder: body.QuestionOrder,
	}

	if err := ctrl.DB.Create(&link).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create question link")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToLectureSessionsQuestionLinkDTO(link))
}
