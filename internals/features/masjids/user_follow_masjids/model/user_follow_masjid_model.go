package model

import (
	"time"
)

type UserFollowMasjid struct {
	FollowUserID    string    `gorm:"type:uuid;not null;primaryKey" json:"follow_user_id"`
	FollowMasjidID  string    `gorm:"type:uuid;not null;primaryKey" json:"follow_masjid_id"`
	FollowCreatedAt time.Time `gorm:"autoCreateTime" json:"follow_created_at"`
}

func (UserFollowMasjid) TableName() string {
	return "user_follow_masjids"
}
