package controller

import (
	"masjidku_backend/internals/features/masjids/lectures/lectures/dto"
	"masjidku_backend/internals/features/masjids/lectures/lectures/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LectureController struct {
	DB *gorm.DB
}

func NewLectureController(db *gorm.DB) *LectureController {
	return &LectureController{DB: db}
}

// 🟢 POST /api/a/lectures
func (ctrl *LectureController) CreateLecture(c *fiber.Ctx) error {
	var req dto.LectureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Permintaan tidak valid", "error": err.Error()})
	}

	newLecture := req.ToModel()
	if err := ctrl.DB.Create(newLecture).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menyimpan data", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Kajian berhasil dibuat",
		"data":    dto.ToLectureResponse(newLecture),
	})
}

// 🟢 GET /api/a/lectures?masjid_id=...
func (ctrl *LectureController) GetLecturesByMasjid(c *fiber.Ctx) error {
	masjidID := c.Query("masjid_id")
	if masjidID == "" {
		return c.Status(400).JSON(fiber.Map{"message": "masjid_id wajib dikirim"})
	}

	var lectures []model.LectureModel
	if err := ctrl.DB.Where("lecture_masjid_id = ?", masjidID).Find(&lectures).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal mengambil data", "error": err.Error()})
	}

	var response []dto.LectureResponse
	for _, l := range lectures {
		response = append(response, *dto.ToLectureResponse(&l))
	}

	return c.JSON(fiber.Map{"message": "Berhasil mengambil daftar kajian", "data": response})
}
