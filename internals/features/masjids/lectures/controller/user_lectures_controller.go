package controller

import (
	"encoding/json"
	"masjidku_backend/internals/features/masjids/lectures/dto"
	"masjidku_backend/internals/features/masjids/lectures/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserLectureController struct {
	DB *gorm.DB
}

func NewUserLectureController(db *gorm.DB) *UserLectureController {
	return &UserLectureController{DB: db}
}

// 🟢 POST /api/a/user-lectures
func (ctrl *UserLectureController) CreateUserLecture(c *fiber.Ctx) error {
	var req dto.UserLectureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Permintaan tidak valid",
			"error":   err.Error(),
		})
	}

	// 🔒 Validasi: pastikan Lecture dan User memang ada
	var count int64
	if err := ctrl.DB.Table("lectures").Where("lecture_id = ?", req.UserLectureLectureID).Count(&count).Error; err != nil || count == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "Lecture tidak ditemukan atau tidak valid"})
	}
	if err := ctrl.DB.Table("users").Where("id = ?", req.UserLectureUserID).Count(&count).Error; err != nil || count == 0 {
		return c.Status(400).JSON(fiber.Map{"message": "User tidak ditemukan atau tidak valid"})
	}

	newUserLecture := req.ToModel()
	if err := ctrl.DB.Create(newUserLecture).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Gagal menyimpan partisipasi",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Partisipasi berhasil dicatat",
		"data":    dto.ToUserLectureResponse(newUserLecture),
	})
}

// 🟢 GET /api/a/user-lectures?lecture_id=...
// 🟢 POST /api/u/user-lectures/by-lecture
func (ctrl *UserLectureController) GetUsersByLecture(c *fiber.Ctx) error {
	// Ambil dari JSON body
	var payload struct {
		LectureID string `json:"lecture_id"`
	}
	if err := c.BodyParser(&payload); err != nil || payload.LectureID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "lecture_id wajib dikirim"})
	}

	// Validasi UUID
	lectureID, err := uuid.Parse(payload.LectureID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "lecture_id tidak valid", "error": err.Error()})
	}

	// Ambil data peserta dari DB
	var participants []model.UserLectureModel
	if err := ctrl.DB.Where("user_lecture_lecture_id = ?", lectureID).Find(&participants).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil peserta", "error": err.Error()})
	}

	// Konversi ke response DTO
	var result []dto.UserLectureResponse
	for _, p := range participants {
		result = append(result, *dto.ToUserLectureResponse(&p))
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil peserta kajian",
		"data":    result,
	})
}

func (ctrl *UserLectureController) GetUserLectureStats(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	userID := ""
	if userIDRaw != nil {
		userID = userIDRaw.(string)
	}

	masjidID := c.Query("masjid_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	month := c.Query("month")
	year := c.Query("year")
	specificDate := c.Query("specific_date")

	if masjidID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Parameter masjid_id wajib diisi")
	}

	type Teacher struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type Result struct {
		model.LectureModel
		UserLectureGradeResult       *int       `json:"user_lecture_grade_result,omitempty"`
		UserLectureCreatedAt         *time.Time `json:"user_lecture_created_at,omitempty"`
		TotalLectureSessions         int        `json:"total_lecture_sessions"`
		CompleteTotalLectureSessions *int       `json:"complete_total_lecture_sessions,omitempty"`
		LectureTeachers              []Teacher  `json:"lecture_teachers,omitempty"`
	}

	// Step 1: Ambil semua data dari tabel lectures (termasuk JSONB pengajar)
	query := ctrl.DB.Table("lectures AS l").
		Select([]string{
			"l.*",
			"l.total_lecture_sessions",
		}).
		Where("l.lecture_masjid_id = ?", masjidID)

	if userID != "" {
		query = query.Select(append(query.Statement.Selects,
			"ul.user_lecture_grade_result",
			"ul.user_lecture_created_at",
			"ul.user_lecture_total_completed_sessions AS complete_total_lecture_sessions",
		)).Joins(`
			LEFT JOIN user_lectures ul 
			ON ul.user_lecture_lecture_id = l.lecture_id 
			AND ul.user_lecture_user_id = ?
		`, userID)
	} else {
		query = query.Select(append(query.Statement.Selects,
			"NULL AS user_lecture_grade_result",
			"NULL AS user_lecture_created_at",
			"NULL AS complete_total_lecture_sessions",
		)).Joins("LEFT JOIN user_lectures ul ON false")
	}

	// Filter waktu
	switch {
	case specificDate != "":
		query = query.Where("DATE(l.lecture_created_at) = ?", specificDate)
	case startDate != "" && endDate != "":
		query = query.Where("l.lecture_created_at BETWEEN ? AND ?", startDate, endDate)
	default:
		if month != "" {
			query = query.Where("EXTRACT(MONTH FROM l.lecture_created_at) = ?", month)
		}
		if year != "" {
			query = query.Where("EXTRACT(YEAR FROM l.lecture_created_at) = ?", year)
		}
	}

	query = query.Order("l.lecture_created_at DESC")

	// Step 2: Eksekusi dan parsing hasil
	type rawLecture struct {
		model.LectureModel
		UserLectureGradeResult       *int       `json:"user_lecture_grade_result"`
		UserLectureCreatedAt         *time.Time `json:"user_lecture_created_at"`
		TotalLectureSessions         int        `json:"total_lecture_sessions"`
		CompleteTotalLectureSessions *int       `json:"complete_total_lecture_sessions"`
	}

	var rawLectures []rawLecture
	if err := query.Scan(&rawLectures).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data kajian")
	}

	// Step 3: Unmarshal JSONB teachers
	var results []Result
	for _, r := range rawLectures {
		var teachers []Teacher
		if r.LectureTeachers != nil {
			_ = json.Unmarshal(r.LectureTeachers, &teachers)
		}

		results = append(results, Result{
			LectureModel:                 r.LectureModel,
			UserLectureGradeResult:       r.UserLectureGradeResult,
			UserLectureCreatedAt:         r.UserLectureCreatedAt,
			TotalLectureSessions:         r.TotalLectureSessions,
			CompleteTotalLectureSessions: r.CompleteTotalLectureSessions,
			LectureTeachers:              teachers,
		})
	}

	// Step 4: Response
	message := "Berhasil mengambil daftar kajian"
	if userID != "" {
		message += " (dengan progress jika login)"
	}

	return c.JSON(fiber.Map{
		"message": message,
		"data":    results,
	})
}
