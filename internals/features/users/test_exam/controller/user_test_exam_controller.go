package controller

import (
	"masjidku_backend/internals/features/users/test_exam/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserTestExamController struct {
	DB *gorm.DB
}

func NewUserTestExamController(db *gorm.DB) *UserTestExamController {
	return &UserTestExamController{DB: db}
}

// ✅ Ambil semua data user_test_exam
func (ctrl *UserTestExamController) GetAll(c *fiber.Ctx) error {
	var results []model.UserTestExam
	if err := ctrl.DB.Order("user_test_exam_id DESC").Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data user_test_exam"})
	}
	return c.JSON(results)
}

// ✅ Ambil satu user_test_exam berdasarkan ID
func (ctrl *UserTestExamController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var data model.UserTestExam
	if err := ctrl.DB.First(&data, "user_test_exam_id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User test exam tidak ditemukan"})
	}
	return c.JSON(data)
}

// ✅ Buat entri user_test_exam baru
func (ctrl *UserTestExamController) Create(c *fiber.Ctx) error {
	var payload model.UserTestExam

	// Parse body JSON ke struct
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Body tidak valid",
		})
	}

	// Ambil user_id dari middleware (disimpan sebagai string)
	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID tidak ditemukan dalam token",
		})
	}

	// Konversi string ke uuid.UUID
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Format User ID tidak valid",
		})
	}

	// Validasi test_exam_id wajib ada
	if payload.UserTestExamTestExamID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Test Exam ID wajib diisi",
		})
	}

	// Set data user_id & timestamp
	payload.UserTestExamUserID = userID
	payload.CreatedAt = time.Now()

	// Simpan ke database
	if err := ctrl.DB.Create(&payload).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan data ke database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User test exam berhasil disimpan",
		"data":    payload,
	})
}

// ✅ Ambil semua hasil test exam berdasarkan user_id
func (ctrl *UserTestExamController) GetByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	var results []model.UserTestExam
	if err := ctrl.DB.
		Where("user_test_exam_user_id = ?", userID).
		Order("created_at DESC").
		Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data berdasarkan user_id",
		})
	}
	return c.JSON(results)
}

// ✅ Ambil semua hasil peserta untuk test_exam tertentu
func (ctrl *UserTestExamController) GetByTestExamID(c *fiber.Ctx) error {
	examID := c.Params("test_exam_id")
	var results []model.UserTestExam
	if err := ctrl.DB.
		Where("user_test_exam_test_exam_id = ?", examID).
		Order("created_at DESC").
		Find(&results).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data berdasarkan test_exam_id",
		})
	}
	return c.JSON(results)
}
