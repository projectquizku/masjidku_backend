package dto

import (
	"masjidku_backend/internals/features/masjids/lecture_sessions/main/model"
	"time"

	"github.com/google/uuid"
)

// ====================
// Response DTO
// ====================

type UserLectureSessionDTO struct {
	UserLectureSessionID               uuid.UUID `json:"user_lecture_session_id"`
	UserLectureSessionAttendanceStatus int       `json:"user_lecture_session_attendance_status"` // 0 = tidak hadir, 1 = hadir, 2 = hadir online
	UserLectureSessionGradeResult      *float64  `json:"user_lecture_session_grade_result,omitempty"`
	UserLectureSessionLectureSessionID uuid.UUID `json:"user_lecture_session_lecture_session_id"`
	UserLectureSessionUserID           uuid.UUID `json:"user_lecture_session_user_id"`
	UserLectureSessionCreatedAt        time.Time `json:"user_lecture_session_created_at"`
}

// ====================
// Request DTO
// ====================

type CreateUserLectureSessionRequest struct {
	UserLectureSessionAttendanceStatus int       `json:"user_lecture_session_attendance_status" validate:"required,oneof=0 1 2"` // 0=tidak hadir, 1=hadir, 2=hadir online
	UserLectureSessionGradeResult      *float64  `json:"user_lecture_session_grade_result,omitempty"`
	UserLectureSessionLectureSessionID uuid.UUID `json:"user_lecture_session_lecture_session_id" validate:"required"`
	UserLectureSessionUserID           uuid.UUID `json:"user_lecture_session_user_id" validate:"required"`
}

// ====================
// Converter
// ====================

func ToUserLectureSessionDTO(u model.UserLectureSessionModel) UserLectureSessionDTO {
	return UserLectureSessionDTO{
		UserLectureSessionID:               u.UserLectureSessionID,
		UserLectureSessionAttendanceStatus: u.UserLectureSessionAttendanceStatus,
		UserLectureSessionGradeResult:      u.UserLectureSessionGradeResult,
		UserLectureSessionLectureSessionID: u.UserLectureSessionLectureSessionID,
		UserLectureSessionUserID:           u.UserLectureSessionUserID,
		UserLectureSessionCreatedAt:        u.UserLectureSessionCreatedAt,
	}
}

func (r CreateUserLectureSessionRequest) ToModel() model.UserLectureSessionModel {
	return model.UserLectureSessionModel{
		UserLectureSessionAttendanceStatus: r.UserLectureSessionAttendanceStatus,
		UserLectureSessionGradeResult:      r.UserLectureSessionGradeResult,
		UserLectureSessionLectureSessionID: r.UserLectureSessionLectureSessionID,
		UserLectureSessionUserID:           r.UserLectureSessionUserID,
	}
}
