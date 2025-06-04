package model

import (
	UserModel "masjidku_backend/internals/features/users/user/model"
	"time"
)

type UserLectureSession struct {
	UserLectureSessionID               string    `gorm:"column:user_lecture_session_id;primaryKey;type:uuid;default:gen_random_uuid()"`
	UserLectureSessionStatusAttendance string    `gorm:"column:user_lecture_session_status_attendance;type:varchar(50)"`
	UserLectureSessionGradeResult      *float64  `gorm:"column:user_lecture_session_grade_result"` // nullable
	UserLectureSessionLectureSessionID string    `gorm:"column:user_lecture_session_lecture_session_id;type:uuid;not null"`
	UserLectureSessionUserID           string    `gorm:"column:user_lecture_session_user_id;type:uuid;not null"`
	UserLectureSessionCreatedAt        time.Time `gorm:"column:user_lecture_session_created_at;autoCreateTime"`

	// Relations
	User           *UserModel.UserModel `gorm:"foreignKey:UserLectureSessionUserID"`
	LectureSession *LectureSessionModel      `gorm:"foreignKey:UserLectureSessionLectureSessionID"`
}

// TableName overrides the default table name
func (UserLectureSession) TableName() string {
	return "user_lecture_sessions"
}
