package controller

import (
	"log"
	"time"

	"masjidku_backend/internals/features/masjids/user_follow_masjids/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserFollowMasjidController struct {
	DB *gorm.DB
}

func NewUserFollowMasjidController(db *gorm.DB) *UserFollowMasjidController {
	return &UserFollowMasjidController{DB: db}
}

// ✅ Follow masjid
func (ctrl *UserFollowMasjidController) FollowMasjid(c *fiber.Ctx) error {
	var input struct {
		MasjidID string `json:"masjid_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format input tidak valid"})
	}

	// Ambil user_id dari JWT claims (via Locals)
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User tidak terautentikasi"})
	}

	userUUID, err1 := uuid.Parse(userIDStr.(string))
	masjidUUID, err2 := uuid.Parse(input.MasjidID)
	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "UUID user atau masjid tidak valid"})
	}

	follow := model.UserFollowMasjidModel{
		FollowUserID:    userUUID,
		FollowMasjidID:  masjidUUID,
		FollowCreatedAt: time.Now(),
	}

	if err := ctrl.DB.Create(&follow).Error; err != nil {
		log.Printf("[ERROR] Gagal follow masjid: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal follow masjid"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil follow masjid",
		"data":    follow,
	})
}

func (ctrl *UserFollowMasjidController) UnfollowMasjid(c *fiber.Ctx) error {
	// Ambil user ID dari JWT token
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User tidak terautentikasi",
		})
	}

	// Ambil masjid ID dari body
	var input struct {
		MasjidID string `json:"masjid_id"`
	}
	if err := c.BodyParser(&input); err != nil || input.MasjidID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Masjid ID harus dikirim dalam body",
		})
	}

	userUUID, err1 := uuid.Parse(userIDStr.(string))
	masjidUUID, err2 := uuid.Parse(input.MasjidID)
	if err1 != nil || err2 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "UUID user atau masjid tidak valid",
		})
	}

	// Delete record follow
	if err := ctrl.DB.Delete(
		&model.UserFollowMasjidModel{},
		"follow_user_id = ? AND follow_masjid_id = ?", userUUID, masjidUUID,
	).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal unfollow masjid",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Berhasil unfollow masjid",
	})
}

// 📄 Lihat semua masjid yang diikuti oleh user (dari JWT token)
func (ctrl *UserFollowMasjidController) GetFollowedMasjidsByUser(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id")
	if userIDStr == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User tidak login"})
	}

	userUUID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	var follows []model.UserFollowMasjidModel
	if err := ctrl.DB.Where("follow_user_id = ?", userUUID).Find(&follows).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal ambil daftar masjid yg di-follow"})
	}

	return c.JSON(fiber.Map{
		"message": "Daftar masjid yang diikuti berhasil diambil",
		"data":    follows,
	})
}
