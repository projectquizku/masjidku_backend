package controller

import (
	"masjidku_backend/internals/features/masjids/lectures/lectures/dto"
	"masjidku_backend/internals/features/masjids/lectures/lectures/model"

	"github.com/gofiber/fiber/v2"
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
func (ctrl *UserLectureController) GetUsersByLecture(c *fiber.Ctx) error {
	lectureID := c.Query("lecture_id")
	if lectureID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "lecture_id wajib dikirim"})
	}

	var participants []model.UserLectureModel
	if err := ctrl.DB.Where("user_lecture_lecture_id = ?", lectureID).Find(&participants).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil peserta", "error": err.Error()})
	}

	var result []dto.UserLectureResponse
	for _, p := range participants {
		result = append(result, *dto.ToUserLectureResponse(&p))
	}

	return c.JSON(fiber.Map{"message": "Berhasil mengambil peserta kajian", "data": result})
}