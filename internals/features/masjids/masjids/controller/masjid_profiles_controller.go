package controller

import (
	"fmt"
	"log"

	"masjidku_backend/internals/features/masjids/masjids/dto"
	"masjidku_backend/internals/features/masjids/masjids/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MasjidProfileController struct {
	DB *gorm.DB
}

func NewMasjidProfileController(db *gorm.DB) *MasjidProfileController {
	return &MasjidProfileController{DB: db}
}

// 🟢 GET PROFILE BY MASJID_ID
func (mpc *MasjidProfileController) GetProfileByMasjidID(c *fiber.Ctx) error {
	masjidIDParam := c.Params("masjid_id")
	log.Printf("[INFO] Fetching profile for masjid ID: %s\n", masjidIDParam)

	// Validasi UUID format
	masjidUUID, err := uuid.Parse(masjidIDParam)
	if err != nil {
		log.Printf("[ERROR] Invalid UUID format: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format UUID masjid tidak valid",
		})
	}

	var profile model.MasjidProfile
	err = mpc.DB.
		Preload("Masjid"). // preload relasi opsional
		Where("masjid_profile_masjid_id = ?", masjidUUID).
		First(&profile).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Printf("[ERROR] Profile not found for masjid ID %s\n", masjidUUID)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Profil masjid tidak ditemukan",
			})
		}

		log.Printf("[ERROR] Database error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Terjadi kesalahan saat mengambil data profil",
		})
	}

	log.Printf("[SUCCESS] Retrieved profile for masjid ID %s\n", masjidUUID)
	return c.JSON(fiber.Map{
		"message": "Profil masjid berhasil diambil",
		"data":    dto.FromModelMasjidProfile(&profile),
	})
}

// 🟢 CREATE PROFILE
func (mpc *MasjidProfileController) CreateMasjidProfile(c *fiber.Ctx) error {
	log.Println("[INFO] Received request to create masjid profile")

	// DTO Input
	var input dto.MasjidProfileRequest
	if err := c.BodyParser(&input); err != nil {
		log.Printf("[ERROR] Invalid input: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format input tidak valid",
		})
	}

	// Validasi Masjid ID (UUID)
	masjidUUID, err := uuid.Parse(input.MasjidProfileMasjidID)
	if err != nil || masjidUUID == uuid.Nil {
		log.Printf("[ERROR] Invalid Masjid ID: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Masjid ID tidak valid atau kosong",
		})
	}

	// Cek apakah sudah ada profile untuk masjid ini
	var existing model.MasjidProfile
	if err := mpc.DB.
		Where("masjid_profile_masjid_id = ?", masjidUUID).
		First(&existing).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Profil untuk masjid ini sudah ada",
		})
	}

	// Konversi ke model
	profile := dto.ToModelMasjidProfile(&input)

	// Simpan ke DB
	if err := mpc.DB.Create(&profile).Error; err != nil {
		log.Printf("[ERROR] Failed to create masjid profile: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan profil masjid",
		})
	}

	log.Printf("[SUCCESS] Masjid profile created for masjid ID: %s\n", masjidUUID)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Profil masjid berhasil dibuat",
		"data":    dto.FromModelMasjidProfile(profile),
	})
}

// 🟢 UPDATE PROFILE (Partial Update)
func (mpc *MasjidProfileController) UpdateMasjidProfile(c *fiber.Ctx) error {
	masjidID := c.Params("masjid_id")
	log.Printf("[INFO] Updating profile for masjid ID: %s\n", masjidID)

	// Ambil data lama dari DB
	var existing model.MasjidProfile
	if err := mpc.DB.Where("masjid_profile_masjid_id = ?", masjidID).First(&existing).Error; err != nil {
		log.Printf("[ERROR] Masjid profile not found: %s\n", masjidID)
		return c.Status(404).JSON(fiber.Map{
			"error": "Profil masjid tidak ditemukan",
		})
	}

	// Bind request ke map untuk partial update
	var inputMap map[string]interface{}
	if err := c.BodyParser(&inputMap); err != nil {
		log.Printf("[ERROR] Invalid input: %v\n", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Format input tidak valid",
		})
	}

	// Hindari update kolom penting
	delete(inputMap, "masjid_profile_id")
	delete(inputMap, "masjid_profile_masjid_id")
	delete(inputMap, "masjid_profile_created_at")

	// Lakukan update
	if err := mpc.DB.Model(&existing).Updates(inputMap).Error; err != nil {
		log.Printf("[ERROR] Failed to update profile: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal memperbarui profil masjid",
		})
	}

	log.Printf("[SUCCESS] Updated profile for masjid ID: %s\n", masjidID)
	return c.JSON(fiber.Map{
		"message": "Profil masjid berhasil diperbarui",
		"data":    dto.FromModelMasjidProfile(&existing),
	})
}

// 🟢 DELETE PROFILE
func (mpc *MasjidProfileController) DeleteMasjidProfile(c *fiber.Ctx) error {
	masjidID := c.Params("masjid_id")
	log.Printf("[INFO] Deleting profile for masjid ID: %s\n", masjidID)

	if err := mpc.DB.Where("masjid_profile_masjid_id = ?", masjidID).
		Delete(&model.MasjidProfile{}).Error; err != nil {
		log.Printf("[ERROR] Failed to delete masjid profile: %v\n", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Gagal menghapus profil masjid",
		})
	}

	log.Printf("[SUCCESS] Masjid profile with masjid ID %s deleted\n", masjidID)
	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Profil masjid dengan ID %s berhasil dihapus", masjidID),
	})
}
