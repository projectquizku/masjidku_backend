package controller

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LectureSessionController struct {
	DB *gorm.DB
}

func NewLectureSessionController(db *gorm.DB) *LectureSessionController {
	return &LectureSessionController{DB: db}
}

// ================================
// CREATE
// ================================
func (ctrl *LectureSessionController) CreateLectureSession(c *fiber.Ctx) error {
	var body dto.CreateLectureSessionRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	// Sementara teacherID dari token, hardcoded dulu jika belum ada JWT
	teacherID := c.Locals("user_id").(string) // ✅ harus disesuaikan jika tidak pakai JWT

	newSession := body.ToModel(teacherID)

	if err := ctrl.DB.Create(&newSession).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create lecture session")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToLectureSessionDTO(newSession))
}

// ================================
// GET ALL
// ================================
func (ctrl *LectureSessionController) GetAllLectureSessions(c *fiber.Ctx) error {
	var sessions []model.LectureSessionModel

	if err := ctrl.DB.Find(&sessions).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to fetch lecture sessions")
	}

	var result []dto.LectureSessionDTO
	for _, s := range sessions {
		result = append(result, dto.ToLectureSessionDTO(s))
	}

	return c.JSON(result)
}

// ================================
// GET BY ID
// ================================
func (ctrl *LectureSessionController) GetLectureSessionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var session model.LectureSessionModel

	if err := ctrl.DB.First(&session, "lecture_session_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Lecture session not found")
	}

	return c.JSON(dto.ToLectureSessionDTO(session))
}

// ✅ POST /api/a/lecture-sessions/by-lecture-id

// ✅ GET lecture sessions by lecture_id (adaptif: jika login, include user progress)
func (ctrl *LectureSessionController) GetByLectureID(c *fiber.Ctx) error {
	type RequestBody struct {
		LectureID string `json:"lecture_id"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil || body.LectureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Permintaan tidak valid, lecture_id wajib diisi",
		})
	}

	// Cek apakah user login
	userID := c.Locals("user_id")

	var sessions []model.LectureSessionModel
	if err := ctrl.DB.
		Where("lecture_session_lecture_id = ?", body.LectureID).
		Order("lecture_session_scheduled_time ASC").
		Find(&sessions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data lecture sessions",
		})
	}

	// Jika tidak login, return langsung
	if userID == nil {
		response := make([]dto.LectureSessionDTO, len(sessions))
		for i, s := range sessions {
			response[i] = dto.ToLectureSessionDTO(s)
		}
		return c.JSON(fiber.Map{
			"message": "Berhasil mengambil lecture sessions",
			"data":    response,
		})
	}

	// Jika login, join ke user_lecture_sessions
	userIDStr := userID.(string)
	type JoinedResult struct {
		model.LectureSessionModel
		AttendanceStatus string   `json:"user_attendance_status"`
		GradeResult      *float64 `json:"user_grade_result"`
	}

	var joined []JoinedResult
	if err := ctrl.DB.Table("lecture_sessions as ls").
		Select(`
				ls.*, 
				uls.user_lecture_session_status_attendance as attendance_status, 
				uls.user_lecture_session_grade_result as grade_result
			`).
		Joins(`LEFT JOIN user_lecture_sessions uls ON uls.user_lecture_session_lecture_session_id = ls.lecture_session_id AND uls.user_lecture_session_user_id = ?`, userIDStr).
		Where("ls.lecture_session_lecture_id = ?", body.LectureID).
		Order("ls.lecture_session_scheduled_time ASC").
		Scan(&joined).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data join lecture sessions",
		})
	}

	// Format responsenya
	response := make([]fiber.Map, len(joined))
	for i, j := range joined {
		resp := dto.ToLectureSessionDTO(j.LectureSessionModel)
		response[i] = fiber.Map{
			"lecture_session":        resp,
			"user_attendance_status": j.AttendanceStatus,
			"user_grade_result":      j.GradeResult,
		}
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil lecture sessions + user progress",
		"data":    response,
	})
}

// ================================
// UPDATE
// ================================
func (ctrl *LectureSessionController) UpdateLectureSession(c *fiber.Ctx) error {
	id := c.Params("id")
	var session model.LectureSessionModel

	if err := ctrl.DB.First(&session, "lecture_session_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Lecture session not found")
	}

	var body dto.UpdateLectureSessionRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	// Update field
	session.LectureSessionTitle = body.LectureSessionTitle
	session.LectureSessionDescription = body.LectureSessionDescription
	session.LectureSessionScheduledTime = body.LectureSessionScheduledTime
	session.LectureSessionPlace = body.LectureSessionPlace
	session.LectureSessionIsSingle = body.LectureSessionIsSingle
	session.LectureSessionLectureID = body.LectureSessionLectureID

	if err := ctrl.DB.Save(&session).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update lecture session")
	}

	return c.JSON(dto.ToLectureSessionDTO(session))
}

// ================================
// DELETE
// ================================
func (ctrl *LectureSessionController) DeleteLectureSession(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.DB.Delete(&model.LectureSessionModel{}, "lecture_session_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete lecture session")
	}

	return c.JSON(fiber.Map{
		"message": "Lecture session deleted successfully",
	})
}
