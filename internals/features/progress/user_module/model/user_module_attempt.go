package model

import (
	"time"

	"github.com/google/uuid"
)

type UserModuleAttempt struct {
	UserModuleAttemptID              uint      `gorm:"primaryKey" json:"user_module_attempt_id"`
	UserModuleAttemptUserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_module_attempt_user_id"`
	UserModuleAttemptTargetType      int       `gorm:"not null" json:"user_module_attempt_target_type"`
	UserModuleAttemptTargetID        int       `gorm:"not null" json:"user_module_attempt_target_id"`
	UserModuleAttemptPercentageGrade *int      `json:"user_module_attempt_percentage_grade"` // bisa nullable
	UserModuleAttemptTimeDuration    *int      `json:"user_module_attempt_time_duration"`    // bisa nullable
	UserModuleAttemptCreatedAt       time.Time `gorm:"autoCreateTime" json:"user_module_attempt_created_at"`

	UserModuleAttemptBatchID     string    `gorm:"type:varchar(100)" json:"user_module_attempt_batch_id"`
	UserModuleAttemptSubmittedAt time.Time `json:"user_module_attempt_submitted_at"`
}

func (UserModuleAttempt) TableName() string {
	return "user_module_attempts"
}
