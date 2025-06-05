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

// 🟢 GET /api/a/lectures/:id
func (ctrl *LectureController) GetLectureByID(c *fiber.Ctx) error {
	lectureID := c.Params("id")
	var lecture model.LectureModel

	if err := ctrl.DB.First(&lecture, "lecture_id = ?", lectureID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Kajian tidak ditemukan", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil mengambil detail kajian",
		"data":    dto.ToLectureResponse(&lecture),
	})
}

// 🟡 PUT /api/a/lectures/:id
func (ctrl *LectureController) UpdateLecture(c *fiber.Ctx) error {
	lectureID := c.Params("id")
	var req dto.LectureRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Permintaan tidak valid", "error": err.Error()})
	}

	var lecture model.LectureModel
	if err := ctrl.DB.First(&lecture, "lecture_id = ?", lectureID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "Kajian tidak ditemukan", "error": err.Error()})
	}

	// Update dengan data baru
	updatedLecture := req.ToModel()
	updatedLecture.LectureID = lecture.LectureID // tetap pakai ID lama

	if err := ctrl.DB.Model(&lecture).Updates(updatedLecture).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal memperbarui data", "error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Kajian berhasil diperbarui",
		"data":    dto.ToLectureResponse(&lecture),
	})
}

// 🔴 DELETE /api/a/lectures/:id
func (ctrl *LectureController) DeleteLecture(c *fiber.Ctx) error {
	lectureID := c.Params("id")

	if err := ctrl.DB.Delete(&model.LectureModel{}, "lecture_id = ?", lectureID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal menghapus kajian", "error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Kajian berhasil dihapus"})
}


