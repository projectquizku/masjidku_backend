package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ThemeModel mewakili tabel themes
type ThemeModel struct {
	ThemeID       uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:theme_id"`
	ThemeName     string         `gorm:"column:theme_name;not null"`
	ThemeType     int            `gorm:"column:theme_type;not null;default:1"`      // 1=color, 2=wallpaper, 3=mix
	ThemeColors   datatypes.JSON `gorm:"type:jsonb;not null;column:theme_colors"`   // JSON berisi warna UI
	Wallpapers    datatypes.JSON `gorm:"type:jsonb;default:'[]';column:wallpapers"` // JSON array wallpaper
	RequiredLevel int            `gorm:"column:required_level;not null"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (ThemeModel) TableName() string {
	return "themes"
}

// Struct pendukung untuk parsing warna dari theme_colors JSONB
type ThemeColors struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Tertiary  string `json:"tertiary"`

	TextPrimary1 string `json:"textPrimary1"`
	TextPrimary2 string `json:"textPrimary2"`
	TextPrimary3 string `json:"textPrimary3"`

	TextSecondary1 string `json:"textSecondary1"`
	TextSecondary2 string `json:"textSecondary2"`
	TextSecondary3 string `json:"textSecondary3"`
	TextSecondary4 string `json:"textSecondary4"`
	TextSecondary5 string `json:"textSecondary5"`

	Success  string `json:"success"`
	Success2 string `json:"success2"`

	Warning  string `json:"warning"`
	Warning1 string `json:"warning1"`

	Accent  string `json:"accent"`
	Accent2 string `json:"accent2"`

	SpecialColor string `json:"specialColor"`
}

// Struct untuk parsing setiap item wallpaper dari JSONB wallpapers
type WallpaperItem struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
	Mode     string `json:"mode"`     // cover / contain / repeat
	Position string `json:"position"` // top / center / bottom
	Tag      string `json:"tag"`      // default / night / day / etc
}

// Fungsi bantu untuk mengambil ThemeColors dari ThemeModel
func (t *ThemeModel) ParseColors() (*ThemeColors, error) {
	var colors ThemeColors
	err := json.Unmarshal(t.ThemeColors, &colors)
	return &colors, err
}

// Fungsi bantu untuk mengambil Wallpapers dari ThemeModel
func (t *ThemeModel) ParseWallpapers() ([]WallpaperItem, error) {
	var wallpapers []WallpaperItem
	err := json.Unmarshal(t.Wallpapers, &wallpapers)
	return wallpapers, err
}
