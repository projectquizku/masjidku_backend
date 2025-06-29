package model

import (
	"time"

	"github.com/google/uuid"
)

type UserLectureModel struct {
	UserLectureID uuid.UUID `gorm:"column:user_lecture_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"user_lecture_id"`

	// Relasi
	UserLectureLectureID uuid.UUID `gorm:"column:user_lecture_lecture_id;type:uuid;not null" json:"user_lecture_lecture_id"`
	UserLectureUserID    uuid.UUID `gorm:"column:user_lecture_user_id;type:uuid;not null" json:"user_lecture_user_id"`

	// Progres & nilai
	UserLectureGradeResult            *int `gorm:"column:user_lecture_grade_result" json:"user_lecture_grade_result,omitempty"`
	UserLectureTotalCompletedSessions int  `gorm:"column:user_lecture_total_completed_sessions;default:0" json:"user_lecture_total_completed_sessions"`

	// Pendaftaran & pembayaran (jika lecture berbayar)
	UserLectureIsRegistered bool       `gorm:"column:user_lecture_is_registered;default:false" json:"user_lecture_is_registered"`
	UserLectureHasPaid      bool       `gorm:"column:user_lecture_has_paid;default:false" json:"user_lecture_has_paid"`
	UserLecturePaidAmount   *int       `gorm:"column:user_lecture_paid_amount" json:"user_lecture_paid_amount,omitempty"`
	UserLecturePaymentTime  *time.Time `gorm:"column:user_lecture_payment_time" json:"user_lecture_payment_time,omitempty"`

	// Metadata
	UserLectureCreatedAt time.Time `gorm:"column:user_lecture_created_at;autoCreateTime" json:"user_lecture_created_at"`
}

// TableName overrides the default table name
func (UserLectureModel) TableName() string {
	return "user_lectures"
}
