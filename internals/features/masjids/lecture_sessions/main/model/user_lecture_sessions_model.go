package model

import (
	UserModel "masjidku_backend/internals/features/users/user/model"
	"time"

	"github.com/google/uuid"
)

type UserLectureSessionModel struct {
	UserLectureSessionID               uuid.UUID `gorm:"column:user_lecture_session_id;primaryKey;type:uuid;default:gen_random_uuid()" json:"user_lecture_session_id"`
	UserLectureSessionAttendanceStatus int       `gorm:"column:user_lecture_session_attendance_status" json:"user_lecture_session_attendance_status"` // 0=tidak hadir, 1=hadir, 2=hadir online
	UserLectureSessionGradeResult      *float64  `gorm:"column:user_lecture_session_grade_result" json:"user_lecture_session_grade_result"`           // nullable
	UserLectureSessionLectureSessionID uuid.UUID `gorm:"column:user_lecture_session_lecture_session_id;type:uuid;not null" json:"user_lecture_session_lecture_session_id"`
	UserLectureSessionUserID           uuid.UUID `gorm:"column:user_lecture_session_user_id;type:uuid;not null" json:"user_lecture_session_user_id"`
	UserLectureSessionCreatedAt        time.Time `gorm:"column:user_lecture_session_created_at;autoCreateTime" json:"user_lecture_session_created_at"`

	// Relations
	User           *UserModel.UserModel `gorm:"foreignKey:UserLectureSessionUserID" json:"user,omitempty"`
	LectureSession *LectureSessionModel `gorm:"foreignKey:UserLectureSessionLectureSessionID" json:"lecture_session,omitempty"`
}

// TableName overrides the default table name
func (UserLectureSessionModel) TableName() string {
	return "user_lecture_sessions"
}
