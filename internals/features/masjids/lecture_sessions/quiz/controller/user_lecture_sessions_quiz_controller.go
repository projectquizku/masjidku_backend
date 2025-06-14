package controller

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/quiz/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/quiz/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserLectureSessionsQuizController struct {
	DB *gorm.DB
}

func NewUserLectureSessionsQuizController(db *gorm.DB) *UserLectureSessionsQuizController {
	return &UserLectureSessionsQuizController{DB: db}
}

// =============================
// ➕ Create User Quiz Result
// =============================
func (ctrl *UserLectureSessionsQuizController) CreateUserLectureSessionsQuiz(c *fiber.Ctx) error {
	var body dto.CreateUserLectureSessionsQuizRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := validate.Struct(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	data := model.UserLectureSessionsQuizModel{
		UserLectureSessionsQuizGrade:  body.UserLectureSessionsQuizGrade,
		UserLectureSessionsQuizQuizID: body.UserLectureSessionsQuizQuizID,
		UserLectureSessionsQuizUserID: body.UserLectureSessionsQuizUserID,
	}

	if err := ctrl.DB.Create(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save quiz result")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToUserLectureSessionsQuizDTO(data))
}

// =============================
// 📄 Get All Quiz Results
// =============================
func (ctrl *UserLectureSessionsQuizController) GetAllUserLectureSessionsQuiz(c *fiber.Ctx) error {
	var results []model.UserLectureSessionsQuizModel

	if err := ctrl.DB.Find(&results).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch quiz results")
	}

	var dtos []dto.UserLectureSessionsQuizDTO
	for _, r := range results {
		dtos = append(dtos, dto.ToUserLectureSessionsQuizDTO(r))
	}

	return c.JSON(dtos)
}

// =============================
// 🔍 Get By Quiz ID or User ID
// =============================
func (ctrl *UserLectureSessionsQuizController) GetUserLectureSessionsQuizFiltered(c *fiber.Ctx) error {
	quizID := c.Query("quiz_id")
	userID := c.Query("user_id")

	query := ctrl.DB.Model(&model.UserLectureSessionsQuizModel{})
	if quizID != "" {
		query = query.Where("user_lecture_sessions_quiz_quiz_id = ?", quizID)
	}
	if userID != "" {
		query = query.Where("user_lecture_sessions_quiz_user_id = ?", userID)
	}

	var results []model.UserLectureSessionsQuizModel
	if err := query.Find(&results).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch filtered quiz results")
	}

	var dtos []dto.UserLectureSessionsQuizDTO
	for _, r := range results {
		dtos = append(dtos, dto.ToUserLectureSessionsQuizDTO(r))
	}

	return c.JSON(dtos)
}

// =============================
// ❌ Delete Quiz Result by ID
// =============================
func (ctrl *UserLectureSessionsQuizController) DeleteUserLectureSessionsQuizByID(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.DB.Delete(&model.UserLectureSessionsQuizModel{}, "user_lecture_sessions_quiz_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete quiz result")
	}

	return c.JSON(fiber.Map{
		"message": "Quiz result deleted successfully",
	})
}
