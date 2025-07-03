package controller

import (
	"log"
	"time"

	"masjidku_backend/internals/features/utils/thema/dto"
	"masjidku_backend/internals/features/utils/thema/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserThemeController struct {
	DB *gorm.DB
}

func NewUserThemeController(db *gorm.DB) *UserThemeController {
	return &UserThemeController{DB: db}
}

// 🟢 Ambil semua tema yang dimiliki user
func (uc *UserThemeController) GetUserThemes(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user_id not found or invalid",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user_id format",
		})
	}

	var userThemes []model.UserThemeModel
	if err := uc.DB.Where("user_id = ?", userID).Find(&userThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user theme"})
	}

	var response []dto.UserThemeResponseDTO
	for _, ut := range userThemes {
		var theme model.ThemeModel
		if err := uc.DB.First(&theme, "theme_id = ?", ut.ThemeID).Error; err != nil {
			continue
		}
		themeDTO, err := dto.MapToThemeResponseDTO(theme)
		if err != nil {
			continue
		}
		response = append(response, dto.MapToUserThemeResponseDTO(ut, &themeDTO))
	}

	return c.JSON(response)
}

// ✅ Ambil tema yang sedang dipilih
func (uc *UserThemeController) GetSelectedTheme(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	var userTheme model.UserThemeModel
	if err := uc.DB.
		Where("user_id = ? AND is_selected = true", userID).
		First(&userTheme).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Belum ada tema yang dipilih"})
	}

	var theme model.ThemeModel
	if err := uc.DB.First(&theme, "theme_id = ?", userTheme.ThemeID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Tema tidak ditemukan"})
	}

	themeDTO, _ := dto.MapToThemeResponseDTO(theme)
	return c.JSON(dto.MapToUserThemeResponseDTO(userTheme, &themeDTO))
}

// 🔓 Unlock tema oleh user
func (uc *UserThemeController) UnlockTheme(c *fiber.Ctx) error {
	// ✅ Amankan parsing user_id
	userIDRaw := c.Locals("user_id")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user_id tidak ditemukan",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id tidak valid",
		})
	}

	// ✅ Parse request body
	var req dto.UserThemeUnlockRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	// ✅ Cek apakah theme_id valid (bisa ditambahkan jika perlu)
	if req.ThemeID == uuid.Nil {
		return c.Status(400).JSON(fiber.Map{"error": "theme_id tidak boleh kosong"})
	}

	// ✅ Cek apakah user sudah unlock theme ini
	var existing model.UserThemeModel
	if err := uc.DB.Where("user_id = ? AND theme_id = ?", userID, req.ThemeID).First(&existing).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tema sudah di-unlock"})
	}

	// ✅ Simpan unlock baru
	newUnlock := model.UserThemeModel{
		UserID:     userID,
		ThemeID:    req.ThemeID,
		IsSelected: false,
		UnlockedAt: time.Now(),
	}

	if err := uc.DB.Create(&newUnlock).Error; err != nil {
		log.Println("[ERROR] Gagal unlock tema:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Gagal unlock tema"})
	}

	return c.JSON(fiber.Map{"message": "Berhasil unlock tema"})
}

// ⭐ Pilih tema sebagai tema aktif
func (uc *UserThemeController) SelectTheme(c *fiber.Ctx) error {
	userIDStr := c.Locals("user_id").(string)
	userID, _ := uuid.Parse(userIDStr)

	var req dto.ThemeSelectRequestDTO
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Request tidak valid"})
	}

	// Pastikan user sudah unlock tema tersebut
	var userTheme model.UserThemeModel
	if err := uc.DB.Where("user_id = ? AND theme_id = ?", userID, req.ThemeID).First(&userTheme).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Tema belum dimiliki"})
	}

	// Reset semua pilihan
	if err := uc.DB.Model(&model.UserThemeModel{}).
		Where("user_id = ?", userID).
		Update("is_selected", false).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal reset pilihan"})
	}

	// Pilih tema baru
	userTheme.IsSelected = true
	userTheme.SelectedWallpaperTag = req.WallpaperTag
	if err := uc.DB.Save(&userTheme).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memilih tema"})
	}

	return c.JSON(fiber.Map{"message": "Berhasil memilih tema"})
}

func (uc *UserThemeController) GetUserThemesWithFallback(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Ambil semua tema
	var allThemes []model.ThemeModel
	if err := uc.DB.Find(&allThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data tema"})
	}

	// Ambil semua user_theme milik user
	var userThemes []model.UserThemeModel
	if err := uc.DB.Where("user_id = ?", userID).Find(&userThemes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data user theme"})
	}

	// Buat map untuk akses cepat
	userThemeMap := make(map[uuid.UUID]model.UserThemeModel)
	for _, ut := range userThemes {
		userThemeMap[ut.ThemeID] = ut
	}

	// Bangun response
	var result []dto.UserThemeResponseDTO
	for _, theme := range allThemes {
		themeDTO, err := dto.MapToThemeResponseDTO(theme)
		if err != nil {
			continue // Skip jika gagal parse JSON string
		}

		if ut, found := userThemeMap[theme.ThemeID]; found {
			result = append(result, dto.MapToUserThemeResponseDTO(ut, &themeDTO))
		} else {
			result = append(result, dto.MapToUserThemeFallback(theme, userID, &themeDTO))
		}
	}

	return c.JSON(result)
}
