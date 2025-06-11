package model

import (
	UserModel "masjidku_backend/internals/features/users/user/model"
	"time"

	"github.com/google/uuid"
)

type LectureSessionModel struct {
	LectureSessionID                     uuid.UUID  `gorm:"column:lecture_session_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"lecture_session_id"`
	LectureSessionTitle                  string     `gorm:"column:lecture_session_title;type:varchar(255);not null" json:"lecture_session_title"`
	LectureSessionDescription            string     `gorm:"column:lecture_session_description;type:text" json:"lecture_session_description"`
	LectureSessionTeacherID              uuid.UUID  `gorm:"column:lecture_session_teacher_id;type:uuid;not null" json:"lecture_session_teacher_id"`
	LectureSessionStartTime              time.Time  `gorm:"column:lecture_session_start_time;not null" json:"lecture_session_start_time"`
	LectureSessionEndTime                time.Time  `gorm:"column:lecture_session_end_time;not null" json:"lecture_session_end_time"`
	LectureSessionPlace                  *string    `gorm:"column:lecture_session_place;type:text" json:"lecture_session_place"`
	LectureSessionLectureID              *uuid.UUID `gorm:"column:lecture_session_lecture_id;type:uuid" json:"lecture_session_lecture_id"`
	LectureSessionMasjidID               uuid.UUID  `gorm:"column:lecture_session_masjid_id;type:uuid;not null" json:"lecture_session_masjid_id"`
	LectureSessionCapacity               *int       `gorm:"column:lecture_session_capacity" json:"lecture_session_capacity"`
	LectureSessionIsPublic               bool       `gorm:"column:lecture_session_is_public;default:true" json:"lecture_session_is_public"`
	LectureSessionIsRegistrationRequired bool       `gorm:"column:lecture_session_is_registration_required;default:false" json:"lecture_session_is_registration_required"`
	LectureSessionIsPaid                 bool       `gorm:"column:lecture_session_is_paid;default:false" json:"lecture_session_is_paid"`
	LectureSessionPrice                  *int       `gorm:"column:lecture_session_price" json:"lecture_session_price"`
	LectureSessionPaymentDeadline        *time.Time `gorm:"column:lecture_session_payment_deadline" json:"lecture_session_payment_deadline"`
	LectureSessionCreatedAt              time.Time  `gorm:"column:lecture_session_created_at;autoCreateTime" json:"lecture_session_created_at"`

	// Relations (optional)
	Teacher *UserModel.UserModel `gorm:"foreignKey:LectureSessionTeacherID" json:"teacher,omitempty"`
	// Lecture *LectureModel         `gorm:"foreignKey:LectureSessionLectureID" json:"lecture,omitempty"` // Uncomment if LectureModel exists
}

// TableName overrides the default table name
func (LectureSessionModel) TableName() string {
	return "lecture_sessions"
}
