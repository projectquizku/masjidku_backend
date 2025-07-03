package model

import (
	"time"

	"github.com/google/uuid"
)

// UserThemeModel merepresentasikan entitas user_themes
type UserThemeModel struct {
	UserThemeID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:user_theme_id"`
	UserID               uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	ThemeID              uuid.UUID `gorm:"type:uuid;not null;column:theme_id"`
	IsSelected           bool      `gorm:"column:is_selected;default:false"`
	SelectedWallpaperTag string    `gorm:"column:selected_wallpaper_tag"`
	UnlockedAt           time.Time `gorm:"column:unlocked_at;autoCreateTime"`
}

// TableName untuk override nama tabel
func (UserThemeModel) TableName() string {
	return "user_themes"
}
