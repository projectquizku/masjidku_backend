package controller

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"masjidku_backend/internals/features/utils/thema/dto"
	"masjidku_backend/internals/features/utils/thema/model"
)

type ThemeController struct {
	DB *gorm.DB
}

func NewThemeController(db *gorm.DB) *ThemeController {
	return &ThemeController{DB: db}
}

// ✅ GET /themes → ambil semua tema
func (tc *ThemeController) GetAllThemes(c *fiber.Ctx) error {
	var themes []model.ThemeModel
	if err := tc.DB.Find(&themes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data tema"})
	}

	var result []dto.ThemeResponseDTO
	for _, theme := range themes {
		dtoData, err := dto.MapToThemeResponseDTO(theme)
		if err != nil {
			continue // skip jika gagal parse
		}
		result = append(result, dtoData)
	}

	return c.JSON(result)
}

// ✅ POST /themes → buat tema baru (admin)
func (tc *ThemeController) CreateTheme(c *fiber.Ctx) error {
	var input dto.ThemeCreateRequestDTO
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Gagal parse data"})
	}

	// 🔁 Konversi warna & wallpaper ke JSON
	themeColorsJSON, err := json.Marshal(input.ThemeColors)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Warna tidak valid"})
	}
	wallpapersJSON, err := json.Marshal(input.Wallpapers)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Wallpapers tidak valid"})
	}

	newTheme := model.ThemeModel{
		ThemeName:     input.ThemeName,
		ThemeType:     input.ThemeType,
		RequiredLevel: input.RequiredLevel,
		ThemeColors:   themeColorsJSON,
		Wallpapers:    wallpapersJSON,
	}

	if err := tc.DB.Create(&newTheme).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan tema"})
	}

	return c.JSON(fiber.Map{
		"message":    "Tema berhasil dibuat",
		"theme_id":   newTheme.ThemeID,
		"theme_name": newTheme.ThemeName,
	})
}
