package model

import (
	UserModel "masjidku_backend/internals/features/users/user/model"
	"time"
)

type LectureSessionModel struct {
	LectureSessionID            string    `gorm:"column:lecture_session_id;primaryKey;type:uuid;default:gen_random_uuid()"`
	LectureSessionTitle         string    `gorm:"column:lecture_session_title;type:varchar(255);not null"`
	LectureSessionDescription   string    `gorm:"column:lecture_session_description;type:text"`
	LectureSessionTeacherID     string    `gorm:"column:lecture_session_teacher_id;type:uuid;not null"`
	LectureSessionScheduledTime time.Time `gorm:"column:lecture_session_scheduled_time;not null"`
	LectureSessionPlace         *string   `gorm:"column:lecture_session_place;type:text"`
	LectureSessionIsSingle      bool      `gorm:"column:lecture_session_is_single;default:false"`
	LectureSessionLectureID     *string   `gorm:"column:lecture_session_lecture_id;type:uuid"`
	LectureSessionCreatedAt     time.Time `gorm:"column:lecture_session_created_at;autoCreateTime"`

	// Relations (optional)
	Teacher *UserModel.UserModel `gorm:"foreignKey:LectureSessionTeacherID"`
	// Lecture *Lecture `gorm:"foreignKey:LectureSessionLectureID"` // Uncomment if Lecture model exists
}

// TableName overrides the default table name
func (LectureSessionModel) TableName() string {
	return "lecture_sessions"
}
