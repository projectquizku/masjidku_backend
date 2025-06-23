package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type LectureModel struct {
	LectureID            uuid.UUID      `gorm:"column:lecture_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"lecture_id"`
	LectureTitle         string         `gorm:"column:lecture_title;type:varchar(255);not null" json:"lecture_title"`
	LectureDescription   string         `gorm:"column:lecture_description;type:text" json:"lecture_description"`
	TotalLectureSessions *int           `gorm:"column:total_lecture_sessions" json:"total_lecture_sessions,omitempty"`
	LectureStatus        bool           `gorm:"column:lecture_status;default:false" json:"lecture_status"` // false = ongoing, true = finished
	LectureCertificateID *uuid.UUID     `gorm:"column:lecture_certificate_id;type:uuid" json:"lecture_certificate_id,omitempty"`
	LectureImageURL      *string        `gorm:"column:lecture_image_url;type:text" json:"lecture_image_url,omitempty"`
	LectureTeachers      datatypes.JSON `gorm:"column:lecture_teachers;type:jsonb" json:"lecture_teachers,omitempty"`
	LectureMasjidID      uuid.UUID      `gorm:"column:lecture_masjid_id;type:uuid;not null" json:"lecture_masjid_id"`
	LectureCreatedAt     time.Time      `gorm:"column:lecture_created_at;autoCreateTime" json:"lecture_created_at"`
}

func (LectureModel) TableName() string {
	return "lectures"
}
