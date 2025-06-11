package controller

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
		return fiber.NewError(fiber.StatusBadRequest, "Permintaan tidak valid")
	}

	// Ambil user_id dari context (harus UUID valid)
	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return fiber.NewError(fiber.StatusUnauthorized, "User ID tidak ditemukan")
	}

	teacherID, err := uuid.Parse(userIDStr)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User ID tidak valid")
	}

	newSession := body.ToModel(teacherID)

	if err := ctrl.DB.Create(&newSession).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal membuat sesi kajian")
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

	// Validasi format UUID
	lectureID, err := uuid.Parse(body.LectureID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Lecture ID tidak valid",
		})
	}

	// Cek apakah user login
	userIDRaw := c.Locals("user_id")

	var sessions []model.LectureSessionModel
	if err := ctrl.DB.
		Where("lecture_session_lecture_id = ?", lectureID).
		Order("lecture_session_start_time ASC").
		Find(&sessions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data sesi kajian",
		})
	}

	// Jika tidak login, return basic response
	if userIDRaw == nil {
		response := make([]dto.LectureSessionDTO, len(sessions))
		for i, s := range sessions {
			response[i] = dto.ToLectureSessionDTO(s)
		}
		return c.JSON(fiber.Map{
			"message": "Berhasil mengambil sesi kajian",
			"data":    response,
		})
	}

	// Parse user_id ke UUID
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User ID tidak valid",
		})
	}

	// LEFT JOIN user progress
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
		Joins(`LEFT JOIN user_lecture_sessions uls 
			ON uls.user_lecture_session_lecture_session_id = ls.lecture_session_id 
			AND uls.user_lecture_session_user_id = ?`, userIDStr).
		Where("ls.lecture_session_lecture_id = ?", lectureID).
		Order("ls.lecture_session_start_time ASC").
		Scan(&joined).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data join sesi kajian",
		})
	}

	// Format response + progress
	response := make([]fiber.Map, len(joined))
	for i, j := range joined {
		sessionDTO := dto.ToLectureSessionDTO(j.LectureSessionModel)
		response[i] = fiber.Map{
			"lecture_session":        sessionDTO,
			"user_attendance_status": j.AttendanceStatus,
			"user_grade_result":      j.GradeResult,
		}
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil sesi kajian + progres user",
		"data":    response,
	})
}

// ================================
// UPDATE
// ================================
func (ctrl *LectureSessionController) UpdateLectureSession(c *fiber.Ctx) error {
	idParam := c.Params("id")

	// Validasi UUID
	sessionID, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID tidak valid")
	}

	var session model.LectureSessionModel
	if err := ctrl.DB.First(&session, "lecture_session_id = ?", sessionID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Sesi kajian tidak ditemukan")
	}

	var body dto.UpdateLectureSessionRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Permintaan tidak valid")
	}

	// Validasi waktu (opsional)
	if body.LectureSessionEndTime.Before(body.LectureSessionStartTime) {
		return fiber.NewError(fiber.StatusBadRequest, "Waktu selesai tidak boleh sebelum waktu mulai")
	}

	// Update semua field
	session.LectureSessionTitle = body.LectureSessionTitle
	session.LectureSessionDescription = body.LectureSessionDescription
	session.LectureSessionStartTime = body.LectureSessionStartTime
	session.LectureSessionEndTime = body.LectureSessionEndTime
	session.LectureSessionPlace = body.LectureSessionPlace
	session.LectureSessionLectureID = body.LectureSessionLectureID
	session.LectureSessionMasjidID = body.LectureSessionMasjidID
	session.LectureSessionCapacity = body.LectureSessionCapacity
	session.LectureSessionIsPublic = body.LectureSessionIsPublic
	session.LectureSessionIsRegistrationRequired = body.LectureSessionIsRegistrationRequired
	session.LectureSessionIsPaid = body.LectureSessionIsPaid
	session.LectureSessionPrice = body.LectureSessionPrice
	session.LectureSessionPaymentDeadline = body.LectureSessionPaymentDeadline

	// Simpan
	if err := ctrl.DB.Save(&session).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui sesi kajian")
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
