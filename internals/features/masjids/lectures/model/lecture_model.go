package model

import (
	"time"

	"github.com/google/uuid"
)

type LectureModel struct {
	LectureID                 uuid.UUID `gorm:"column:lecture_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"lecture_id"`
	LectureTitle              string    `gorm:"column:lecture_title;type:varchar(255);not null" json:"lecture_title"`
	LectureDescription        string    `gorm:"column:lecture_description;type:text" json:"lecture_description"`
	TotalLectureSessions      *int      `gorm:"column:total_lecture_sessions" json:"total_lecture_sessions"` // nullable
	LectureIsRecurring        bool      `gorm:"column:lecture_is_recurring;default:false" json:"lecture_is_recurring"`
	LectureRecurrenceInterval *int      `gorm:"column:lecture_recurrence_interval" json:"lecture_recurrence_interval"` // nullable
	LectureMasjidID           uuid.UUID `gorm:"column:lecture_masjid_id;type:uuid;not null" json:"lecture_masjid_id"`
	LectureCreatedAt          time.Time `gorm:"column:lecture_created_at;autoCreateTime" json:"lecture_created_at"`
}

// TableName overrides the table name
func (LectureModel) TableName() string {
	return "lectures"
}
