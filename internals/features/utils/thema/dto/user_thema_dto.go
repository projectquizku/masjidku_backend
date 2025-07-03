package dto

import (
	"time"

	"masjidku_backend/internals/features/utils/thema/model"

	"github.com/google/uuid"
)

// ✅ DTO untuk menampilkan tema yang dimiliki user
type UserThemeResponseDTO struct {
	UserThemeID          uuid.UUID         `json:"user_theme_id"`
	UserID               uuid.UUID         `json:"user_id"`
	ThemeID              uuid.UUID         `json:"theme_id"`
	IsSelected           bool              `json:"is_selected"`
	SelectedWallpaperTag string            `json:"selected_wallpaper_tag"`
	UnlockedAt           *time.Time        `json:"unlocked_at"` // ✅ pakai pointer agar bisa null
	Theme                *ThemeResponseDTO `json:"theme"`
	IsUnlocked           bool              `json:"is_unlocked"`
}

// ✅ DTO untuk unlock tema (misal saat user mencapai level tertentu)
type UserThemeUnlockRequestDTO struct {
	ThemeID uuid.UUID `json:"theme_id" validate:"required"`
}

// ✅ DTO untuk memilih tema aktif (bisa disatukan dengan yang sudah ada)
type ThemeSelectRequestDTO struct {
	ThemeID      uuid.UUID `json:"theme_id" validate:"required"`
	WallpaperTag string    `json:"wallpaper_tag,omitempty"` // opsional
}

// ✅ Fungsi bantu untuk memetakan model ke response
func MapToUserThemeResponseDTO(userTheme model.UserThemeModel, theme *ThemeResponseDTO) UserThemeResponseDTO {
	return UserThemeResponseDTO{
		UserThemeID:          userTheme.UserThemeID,
		UserID:               userTheme.UserID,
		ThemeID:              userTheme.ThemeID,
		IsSelected:           userTheme.IsSelected,
		SelectedWallpaperTag: userTheme.SelectedWallpaperTag,
		UnlockedAt:           &userTheme.UnlockedAt, // ✅ pointer
		Theme:                theme,
		IsUnlocked:           true, // ✅ eksplisit unlocked
	}
}

// ✅ Fungsi fallback jika user belum punya tema tersebut
func MapToUserThemeFallback(theme model.ThemeModel, userID uuid.UUID, themeDTO *ThemeResponseDTO) UserThemeResponseDTO {
	return UserThemeResponseDTO{
		UserThemeID:          uuid.Nil,
		UserID:               userID,
		ThemeID:              theme.ThemeID,
		IsSelected:           false,
		SelectedWallpaperTag: "",
		UnlockedAt:           nil, // ✅ null
		Theme:                themeDTO,
		IsUnlocked:           false,
	}
}
