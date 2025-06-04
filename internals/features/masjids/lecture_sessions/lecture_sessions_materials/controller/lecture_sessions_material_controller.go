package controller

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_materials/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions_materials/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate2 = validator.New() // ✅ Buat instance validator

type LectureSessionsMaterialController struct {
	DB *gorm.DB
}

func NewLectureSessionsMaterialController(db *gorm.DB) *LectureSessionsMaterialController {
	return &LectureSessionsMaterialController{DB: db}
}

// =============================
// ➕ Create Lecture Session Material
// =============================
func (ctrl *LectureSessionsMaterialController) CreateLectureSessionsMaterial(c *fiber.Ctx) error {
	var body dto.CreateLectureSessionsMaterialRequest
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := validate2.Struct(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	material := model.LectureSessionsMaterialModel{
		LectureSessionsMaterialTitle:            body.LectureSessionsMaterialTitle,
		LectureSessionsMaterialSummary:          body.LectureSessionsMaterialSummary,
		LectureSessionsMaterialTranscriptFull:   body.LectureSessionsMaterialTranscriptFull,
		LectureSessionsMaterialLectureSessionID: body.LectureSessionsMaterialLectureSessionID,
	}

	if err := ctrl.DB.Create(&material).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create material")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToLectureSessionsMaterialDTO(material))
}

// =============================
// 📄 Get All Materials
// =============================
func (ctrl *LectureSessionsMaterialController) GetAllLectureSessionsMaterials(c *fiber.Ctx) error {
	var materials []model.LectureSessionsMaterialModel

	if err := ctrl.DB.Find(&materials).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve materials")
	}

	var response []dto.LectureSessionsMaterialDTO
	for _, m := range materials {
		response = append(response, dto.ToLectureSessionsMaterialDTO(m))
	}

	return c.JSON(response)
}

// =============================
// 🔍 Get Material by ID
// =============================
func (ctrl *LectureSessionsMaterialController) GetLectureSessionsMaterialByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var material model.LectureSessionsMaterialModel
	if err := ctrl.DB.First(&material, "lecture_sessions_material_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Material not found")
	}

	return c.JSON(dto.ToLectureSessionsMaterialDTO(material))
}

// =============================
// ❌ Delete Material
// =============================
func (ctrl *LectureSessionsMaterialController) DeleteLectureSessionsMaterial(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.DB.Delete(&model.LectureSessionsMaterialModel{}, "lecture_sessions_material_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete material")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
