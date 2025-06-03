package model

import (
	"time"

	"github.com/google/uuid"
)

type UserFollowMasjid struct {
	FollowUserID    uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"follow_user_id"`
	FollowMasjidID  uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"follow_masjid_id"`
	FollowCreatedAt time.Time `gorm:"autoCreateTime" json:"follow_created_at"`
}

func (UserFollowMasjid) TableName() string {
	return "user_follow_masjid"
}
