package controller

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/quiz/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/quiz/model"
	"time"

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
// ➕ Create User Quiz Result (from token)
// =============================
func (ctrl *UserLectureSessionsQuizController) CreateUserLectureSessionsQuiz(c *fiber.Ctx) error {
	var body dto.CreateUserLectureSessionsQuizRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if err := validate.Struct(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Ambil user_id dari token (diset oleh middleware)
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User ID not found in token")
	}
	userID := userIDRaw.(string)

	data := model.UserLectureSessionsQuizModel{
		UserLectureSessionsQuizGrade:  body.UserLectureSessionsQuizGrade,
		UserLectureSessionsQuizQuizID: body.UserLectureSessionsQuizQuizID,
		UserLectureSessionsQuizUserID: userID,
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



func (ctrl *UserLectureSessionsQuizController) GetUserQuizWithDetail(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User ID not found in token")
	}
	userID := userIDRaw.(string)

	lectureID := c.Query("lecture_id")
	lectureSessionID := c.Query("lecture_session_id")

	type UserQuizWithDetail struct {
		UserLectureSessionsQuizID        string    `json:"user_lecture_sessions_quiz_id"`
		UserLectureSessionsQuizGrade     float64   `json:"user_lecture_sessions_quiz_grade_result"`
		UserLectureSessionsQuizUserID    string    `json:"user_lecture_sessions_quiz_user_id"`
		UserLectureSessionsQuizCreatedAt time.Time `json:"user_lecture_sessions_quiz_created_at"`

		LectureSessionsQuizID               string    `json:"lecture_sessions_quiz_id"`
		LectureSessionsQuizTitle            string    `json:"lecture_sessions_quiz_title"`
		LectureSessionsQuizDescription      string    `json:"lecture_sessions_quiz_description"`
		LectureSessionsQuizLectureSessionID string    `json:"lecture_sessions_quiz_lecture_session_id"`
		LectureSessionsQuizCreatedAt        time.Time `json:"lecture_sessions_quiz_created_at"`
	}

	var results []UserQuizWithDetail

	query := ctrl.DB.
		Table("user_lecture_sessions_quiz AS uq").
		Select(`
			uq.user_lecture_sessions_quiz_id,
			uq.user_lecture_sessions_quiz_grade_result,
			uq.user_lecture_sessions_quiz_user_id,
			uq.user_lecture_sessions_quiz_created_at,

			q.lecture_sessions_quiz_id,
			q.lecture_sessions_quiz_title,
			q.lecture_sessions_quiz_description,
			q.lecture_sessions_quiz_lecture_session_id,
			q.lecture_sessions_quiz_created_at
		`).
		Joins("JOIN lecture_sessions_quiz AS q ON uq.user_lecture_sessions_quiz_quiz_id = q.lecture_sessions_quiz_id").
		Where("uq.user_lecture_sessions_quiz_user_id = ?", userID)

	if lectureID != "" {
		query = query.Joins("JOIN lecture_sessions AS ls ON q.lecture_sessions_quiz_lecture_session_id = ls.lecture_session_id").
			Where("ls.lecture_session_lecture_id = ?", lectureID)
	}

	if lectureSessionID != "" {
		query = query.Where("q.lecture_sessions_quiz_lecture_session_id = ?", lectureSessionID)
	}

	err := query.Scan(&results).Error
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch quiz results with details")
	}

	return c.JSON(results)
}
