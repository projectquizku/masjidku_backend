package dto

import (
	"masjidku_backend/internals/features/utils/thema/model"

	"github.com/google/uuid"
)

// ✅ DTO untuk response data tema ke frontend
type ThemeResponseDTO struct {
	ThemeID       uuid.UUID             `json:"theme_id"`
	ThemeName     string                `json:"theme_name"`
	ThemeType     int                   `json:"theme_type"` // 1=color, 2=wallpaper, 3=mix
	RequiredLevel int                   `json:"required_level"`
	ThemeColors   model.ThemeColors     `json:"theme_colors"`
	Wallpapers    []model.WallpaperItem `json:"wallpapers"`
}

// ✅ DTO untuk membuat tema baru (admin)
type ThemeCreateRequestDTO struct {
	ThemeName     string                `json:"theme_name" validate:"required"`
	ThemeType     int                   `json:"theme_type" validate:"required,oneof=1 2 3"`
	RequiredLevel int                   `json:"required_level" validate:"required,min=1"`
	ThemeColors   model.ThemeColors     `json:"theme_colors" validate:"required"`
	Wallpapers    []model.WallpaperItem `json:"wallpapers"`
}




// ✅ Fungsi bantu mapping model ke response DTO
func MapToThemeResponseDTO(theme model.ThemeModel) (ThemeResponseDTO, error) {
	colors, err := theme.ParseColors()
	if err != nil {
		return ThemeResponseDTO{}, err
	}

	wallpapers, err := theme.ParseWallpapers()
	if err != nil {
		return ThemeResponseDTO{}, err
	}

	return ThemeResponseDTO{
		ThemeID:       theme.ThemeID,
		ThemeName:     theme.ThemeName,
		ThemeType:     theme.ThemeType,
		RequiredLevel: theme.RequiredLevel,
		ThemeColors:   *colors,
		Wallpapers:    wallpapers,
	}, nil
}
