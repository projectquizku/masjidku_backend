package controller

import (
	"log"
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"
	helper "masjidku_backend/internals/helpers"
	"time"

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

func (ctrl *LectureSessionController) CreateLectureSession(c *fiber.Ctx) error {
	log.Println("[INFO] Menerima request untuk membuat sesi kajian")

	// Validasi user login (optional)
	userIDRaw := c.Locals("user_id")
	if userIDRaw == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "User belum login")
	}

	// Ambil dan validasi field wajib
	title := c.FormValue("lecture_session_title")
	startStr := c.FormValue("lecture_session_start_time")
	endStr := c.FormValue("lecture_session_end_time")
	teacherJSON := c.FormValue("lecture_session_teacher")

	if title == "" || startStr == "" || endStr == "" || teacherJSON == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Field wajib tidak lengkap")
	}

	// Parse waktu
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format start time tidak valid (RFC3339)")
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format end time tidak valid (RFC3339)")
	}

	// Parse teacher JSON
	var teacher dto.JSONBTeacher
	if err := teacher.FromString(teacherJSON); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Format teacher tidak valid")
	}

	// Parse opsional UUID
	var lectureID *uuid.UUID
	if val := c.FormValue("lecture_session_lecture_id"); val != "" {
		if parsed, err := uuid.Parse(val); err == nil {
			lectureID = &parsed
		}
	}

	var certificateID *uuid.UUID
	if val := c.FormValue("lecture_session_certificate_id"); val != "" {
		if parsed, err := uuid.Parse(val); err == nil {
			certificateID = &parsed
		}
	}

	// Parse opsional string dan waktu
	description := c.FormValue("lecture_session_description")
	place := c.FormValue("lecture_session_place")

	// Upload gambar jika ada
	var imageURL *string
	if file, err := c.FormFile("lecture_session_image_url"); err == nil && file != nil {
		url, err := helper.UploadImageAsWebPToSupabase("lecture_sessions", file)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal upload gambar")
		}
		imageURL = &url
	} else if val := c.FormValue("lecture_session_image_url"); val != "" {
		imageURL = &val
	}

	// Susun DTO
	req := dto.CreateLectureSessionRequest{
		LectureSessionTitle:         title,
		LectureSessionDescription:   description,
		LectureSessionTeacher:       teacher,
		LectureSessionStartTime:     startTime,
		LectureSessionEndTime:       endTime,
		LectureSessionPlace:         &place,
		LectureSessionImageURL:      imageURL,
		LectureSessionLectureID:     lectureID,
		LectureSessionCertificateID: certificateID,
	}

	// Simpan ke DB
	lectureSession := req.ToModel()
	if err := ctrl.DB.Create(&lectureSession).Error; err != nil {
		log.Println("[ERROR] Gagal menyimpan ke DB:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan sesi kajian")
	}

	// Response sukses
	return c.Status(fiber.StatusCreated).JSON(dto.ToLectureSessionDTO(lectureSession))
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

	lectureID, err := uuid.Parse(body.LectureID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Lecture ID tidak valid",
		})
	}

	userIDRaw := c.Locals("user_id")

	// Jika tidak login, ambil data biasa
	if userIDRaw == nil {
		var sessions []model.LectureSessionModel
		if err := ctrl.DB.
			Where("lecture_session_lecture_id = ?", lectureID).
			Order("lecture_session_start_time ASC").
			Find(&sessions).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Gagal mengambil data sesi kajian",
			})
		}

		response := make([]dto.LectureSessionDTO, len(sessions))
		for i, s := range sessions {
			response[i] = dto.ToLectureSessionDTO(s)
		}

		return c.JSON(fiber.Map{
			"message": "Berhasil mengambil sesi kajian",
			"data":    response,
		})
	}

	// Jika login → Ambil juga progress user
	userIDStr, ok := userIDRaw.(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User ID tidak valid",
		})
	}

	// Ambil data + progress via LEFT JOIN
	type JoinedResult struct {
		model.LectureSessionModel
		UserAttendanceStatus string   `json:"user_attendance_status"`
		UserGradeResult      *float64 `json:"user_grade_result"`
	}

	var joined []JoinedResult
	if err := ctrl.DB.Table("lecture_sessions as ls").
		Select(`
			ls.*, 
			uls.user_lecture_session_status_attendance as user_attendance_status, 
			uls.user_lecture_session_grade_result as user_grade_result
		`).
		Joins(`
			LEFT JOIN user_lecture_sessions uls 
			ON uls.user_lecture_session_lecture_session_id = ls.lecture_session_id 
			AND uls.user_lecture_session_user_id = ?
		`, userIDStr).
		Where("ls.lecture_session_lecture_id = ?", lectureID).
		Order("ls.lecture_session_start_time ASC").
		Scan(&joined).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data sesi + progres user",
		})
	}

	// Gabungkan ke response
	response := make([]fiber.Map, len(joined))
	for i, j := range joined {
		response[i] = fiber.Map{
			"lecture_session":        dto.ToLectureSessionDTO(j.LectureSessionModel),
			"user_attendance_status": j.UserAttendanceStatus,
			"user_grade_result":      j.UserGradeResult,
		}
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil sesi kajian + progres user",
		"data":    response,
	})
}

func (ctrl *LectureSessionController) UpdateLectureSession(c *fiber.Ctx) error {
	idParam := c.Params("id")
	sessionID, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID tidak valid")
	}

	// Ambil data existing
	var existing model.LectureSessionModel
	if err := ctrl.DB.First(&existing, "lecture_session_id = ?", sessionID).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Sesi kajian tidak ditemukan")
	}

	// Bind JSON ke struct sementara
	var body dto.UpdateLectureSessionRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Permintaan tidak valid")
	}

	// Validasi waktu
	if !body.LectureSessionStartTime.IsZero() && !body.LectureSessionEndTime.IsZero() {
		if body.LectureSessionEndTime.Before(body.LectureSessionStartTime) {
			return fiber.NewError(fiber.StatusBadRequest, "Waktu selesai tidak boleh sebelum waktu mulai")
		}
	}

	// ✅ Upload gambar baru jika ada
	fileHeader, err := c.FormFile("lecture_session_image_url")
	if err == nil && fileHeader != nil {
		// 🔥 Hapus gambar lama
		if existing.LectureSessionImageURL != nil {
			_ = helper.DeleteFromSupabase("image", helper.ExtractSupabaseStoragePath(*existing.LectureSessionImageURL))
		}

		// ⬆️ Upload gambar baru
		url, err := helper.UploadImageAsWebPToSupabase("lecture_sessions", fileHeader)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Gagal upload gambar baru")
		}
		existing.LectureSessionImageURL = &url
	}

	// ✅ Update field satu per satu jika ada isinya
	if body.LectureSessionTitle != "" {
		existing.LectureSessionTitle = body.LectureSessionTitle
	}
	if body.LectureSessionDescription != "" {
		existing.LectureSessionDescription = body.LectureSessionDescription
	}
	if body.LectureSessionTeacher.ID != "" && body.LectureSessionTeacher.Name != "" {
		existing.LectureSessionTeacher = body.LectureSessionTeacher.ToModel()
	}
	if !body.LectureSessionStartTime.IsZero() {
		existing.LectureSessionStartTime = body.LectureSessionStartTime
	}
	if !body.LectureSessionEndTime.IsZero() {
		existing.LectureSessionEndTime = body.LectureSessionEndTime
	}
	if body.LectureSessionPlace != nil {
		existing.LectureSessionPlace = body.LectureSessionPlace
	}
	if body.LectureSessionLectureID != nil {
		existing.LectureSessionLectureID = body.LectureSessionLectureID
	}
	if body.LectureSessionCertificateID != nil {
		existing.LectureSessionCertificateID = body.LectureSessionCertificateID
	}

	// Simpan ke database
	if err := ctrl.DB.Save(&existing).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui sesi kajian")
	}

	// ✅ Kirim response dalam format DTO
	return c.JSON(dto.ToLectureSessionDTO(existing))
}

// ================================
// DELETE
// ================================
func (ctrl *LectureSessionController) DeleteLectureSession(c *fiber.Ctx) error {
	idParam := c.Params("id")
	sessionID, err := uuid.Parse(idParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "ID tidak valid")
	}

	if err := ctrl.DB.Delete(&model.LectureSessionModel{}, "lecture_session_id = ?", sessionID).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menghapus sesi kajian")
	}

	return c.JSON(fiber.Map{
		"message": "Sesi kajian berhasil dihapus",
	})
}
