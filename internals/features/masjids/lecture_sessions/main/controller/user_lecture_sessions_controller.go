package controller

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/dto"
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserLectureSessionController struct {
	DB *gorm.DB
}

func NewUserLectureSessionController(db *gorm.DB) *UserLectureSessionController {
	return &UserLectureSessionController{DB: db}
}

// CREATE
func (ctrl *UserLectureSessionController) CreateUserLectureSession(c *fiber.Ctx) error {
	var req dto.CreateUserLectureSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Permintaan tidak valid")
	}

	newRecord := model.UserLectureSessionModel{
		UserLectureSessionAttendanceStatus: req.UserLectureSessionAttendanceStatus,
		UserLectureSessionGradeResult:      req.UserLectureSessionGradeResult,
		UserLectureSessionLectureSessionID: req.UserLectureSessionLectureSessionID,
		UserLectureSessionUserID:           req.UserLectureSessionUserID,
		UserLectureSessionIsRegistered:     req.UserLectureSessionIsRegistered,
		UserLectureSessionHasPaid:          req.UserLectureSessionHasPaid,
		UserLectureSessionPaidAmount:       req.UserLectureSessionPaidAmount,
		UserLectureSessionPaymentTime:      req.UserLectureSessionPaymentTime,
	}

	if err := ctrl.DB.Create(&newRecord).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal membuat user lecture session")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToUserLectureSessionDTO(newRecord))
}

// GET ALL
func (ctrl *UserLectureSessionController) GetAllUserLectureSessions(c *fiber.Ctx) error {
	var records []model.UserLectureSessionModel
	if err := ctrl.DB.Find(&records).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve records")
	}

	var result []dto.UserLectureSessionDTO
	for _, record := range records {
		result = append(result, dto.ToUserLectureSessionDTO(record))
	}

	return c.JSON(result)
}

// GET BY ID
func (ctrl *UserLectureSessionController) GetUserLectureSessionByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var record model.UserLectureSessionModel
	if err := ctrl.DB.First(&record, "user_lecture_session_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Record not found")
	}

	return c.JSON(dto.ToUserLectureSessionDTO(record))
}

// UPDATE
func (ctrl *UserLectureSessionController) UpdateUserLectureSession(c *fiber.Ctx) error {
	id := c.Params("id")

	var record model.UserLectureSessionModel
	if err := ctrl.DB.First(&record, "user_lecture_session_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Data tidak ditemukan")
	}

	var req dto.CreateUserLectureSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Permintaan tidak valid")
	}

	// Update field
	record.UserLectureSessionAttendanceStatus = req.UserLectureSessionAttendanceStatus
	record.UserLectureSessionGradeResult = req.UserLectureSessionGradeResult
	record.UserLectureSessionLectureSessionID = req.UserLectureSessionLectureSessionID
	record.UserLectureSessionUserID = req.UserLectureSessionUserID
	record.UserLectureSessionIsRegistered = req.UserLectureSessionIsRegistered
	record.UserLectureSessionHasPaid = req.UserLectureSessionHasPaid
	record.UserLectureSessionPaidAmount = req.UserLectureSessionPaidAmount
	record.UserLectureSessionPaymentTime = req.UserLectureSessionPaymentTime

	if err := ctrl.DB.Save(&record).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal memperbarui data user lecture session")
	}

	return c.JSON(dto.ToUserLectureSessionDTO(record))
}

// DELETE
func (ctrl *UserLectureSessionController) DeleteUserLectureSession(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := ctrl.DB.Delete(&model.UserLectureSessionModel{}, "user_lecture_session_id = ?", id).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete")
	}

	return c.SendStatus(fiber.StatusNoContent)
}
