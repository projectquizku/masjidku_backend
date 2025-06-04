package dto

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/lecture_sessions/model"
	"time"
)

// ====================
// Response DTO
// ====================

type UserLectureSessionDTO struct {
	UserLectureSessionID               string    `json:"user_lecture_session_id"`
	UserLectureSessionStatusAttendance string    `json:"user_lecture_session_status_attendance"`
	UserLectureSessionGradeResult      *float64  `json:"user_lecture_session_grade_result,omitempty"`
	UserLectureSessionLectureSessionID string    `json:"user_lecture_session_lecture_session_id"`
	UserLectureSessionUserID           string    `json:"user_lecture_session_user_id"`
	UserLectureSessionCreatedAt        time.Time `json:"user_lecture_session_created_at"`
}

// ====================
// Request DTO
// ====================

type CreateUserLectureSessionRequest struct {
	UserLectureSessionStatusAttendance string   `json:"user_lecture_session_status_attendance" validate:"required,oneof=hadir tidak_hadir izin"`
	UserLectureSessionGradeResult      *float64 `json:"user_lecture_session_grade_result,omitempty"`
	UserLectureSessionLectureSessionID string   `json:"user_lecture_session_lecture_session_id" validate:"required,uuid"`
	UserLectureSessionUserID           string   `json:"user_lecture_session_user_id" validate:"required,uuid"`
}

// ====================
// Converter
// ====================

func ToUserLectureSessionDTO(u model.UserLectureSession) UserLectureSessionDTO {
	return UserLectureSessionDTO{
		UserLectureSessionID:               u.UserLectureSessionID,
		UserLectureSessionStatusAttendance: u.UserLectureSessionStatusAttendance,
		UserLectureSessionGradeResult:      u.UserLectureSessionGradeResult,
		UserLectureSessionLectureSessionID: u.UserLectureSessionLectureSessionID,
		UserLectureSessionUserID:           u.UserLectureSessionUserID,
		UserLectureSessionCreatedAt:        u.UserLectureSessionCreatedAt,
	}
}
